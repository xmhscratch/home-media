package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"home-media/sys"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (ctx *SQSegmentInfo) Init(cfg *sys.Config) error {
	var (
		err            error
		best           [3]float64           = [3]float64{0, 0, 0}
		rds            *redis.Client        = sys.NewClient(cfg)
		genSegmentPath func(c int64) string = func(c int64) string {
			return filepath.Join(
				cfg.DataPath, ctx.SavePath[:24],
				fmt.Sprintf("%s_%03d", GetFileKeyName(ctx.SavePath), c),
			)
		}
	)
	ctx.Config = cfg
	defer rds.Close()

	ctx.TotalLength = ctx.FileMeta.Duration
	// fmt.Println(ctx.TotalLength, ctx.FileMeta.Duration)
	best = ctx.bestSegmentValue()
	ctx.BestSegmentLength = (time.Duration(best[0]) * time.Minute).Seconds()
	ctx.BestSegmentCount = int64(best[1])

	// litter.D(best)

	ctx.Segments = map[string]string{}
	for c := range ctx.BestSegmentCount {
		segmentPath := genSegmentPath(c)
		startDuration := float64(c) * ctx.BestSegmentLength
		ctx.Segments[segmentPath] = FormatDuration(startDuration)
	}

	if err == nil {
		ctx.KeyId = GetFileKeyName(ctx.SavePath)

		if err := rds.HSet(
			context.TODO(),
			GetKeyName("segment", ":count"),
			[]string{ctx.KeyId, strconv.FormatInt(ctx.BestSegmentCount, 5<<1)},
		).Err(); err != nil {
			return err
		}
	}

	return err
}

func (ctx *SQSegmentInfo) PushQueue() error {
	var (
		// err error
		rds    *redis.Client = sys.NewClient(ctx.Config)
		c      float64       = 0.0000
		zitems []redis.Z     = []redis.Z{}
	)
	defer rds.Close()

	if ctx.KeyId == "" {
		return errors.New("queue is not initialized")
	}

	for output, start := range ctx.Segments {
		zitems = append(zitems, redis.Z{
			Score: c,
			Member: &SQMessage{
				KeyId:    ctx.KeyId,
				Index:    int64(c),
				Source:   filepath.Join(ctx.Config.DataPath, ctx.SavePath),
				Start:    start,
				Duration: FormatDuration(ctx.BestSegmentLength),
				Output:   output,
			},
		})
		c += 1
	}

	for output := range ctx.Segments {
		if err := rds.SAdd(
			sys.SessionContext,
			GetKeyName("concat:queue", ":", ctx.KeyId),
			output,
		).Err(); err != nil {
			return err
		}
	}

	return rds.ZAddNX(
		sys.SessionContext,
		GetKeyName("segment", ":queue"),
		zitems...,
	).Err()
}

func (ctx *SQSegmentInfo) bestSegmentValue() [3]float64 {
	var (
		best        chan [3]float64           = make(chan [3]float64, 1)
		pbSegLength [SEGMENT_CAPACITY]float64 = [SEGMENT_CAPACITY]float64{2, 3, 5, 7, 11, 13}
		pbSegCount  []float64                 = []float64{}
		fnGetBest   func(i int) [3]float64    = func(i int) [3]float64 {
			pbSegLengthVal := pbSegLength[i]
			pbSegCountVal := math.Ceil(ctx.TotalLength / minToSec(pbSegLengthVal))
			pbSegPeriodVal := math.Abs(pbSegCountVal - pbSegLengthVal)

			return [3]float64{pbSegLengthVal, pbSegCountVal, pbSegPeriodVal}
		}
	)
	defer close(best)

	best <- fnGetBest(0)

findLoop:
	for {
		if stop := func(i chan int) bool {
			i <- len(pbSegCount)

			newBest := fnGetBest(<-i)
			pbSegLengthVal := newBest[0]
			pbSegCountVal := newBest[1]
			pbSegPeriodVal := newBest[2]

			pbSegCount = append(pbSegCount, pbSegCountVal)

			var currBest [3]float64 = <-best
			if pbSegPeriodVal < currBest[2] {
				best <- [3]float64{pbSegLengthVal, pbSegCountVal, pbSegPeriodVal}
			} else {
				best <- currBest
			}

			return len(pbSegCount) == len(pbSegLength)
		}(make(chan int, 1)); stop {
			break findLoop
		}
	}

	return <-best
}

func (ctx SQMessage) MarshalBinary() (data []byte, err error) {
	return json.Marshal(ctx)
}
