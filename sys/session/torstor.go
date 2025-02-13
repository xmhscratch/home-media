package session

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/anacrolix/log"
	"github.com/anacrolix/missinggo/v2"

	"github.com/anacrolix/torrent/common"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/segments"
	"github.com/anacrolix/torrent/storage"
)

// File-based storage for torrents, that isn't yet bound to a particular torrent.
type fileClientImpl struct {
	opts NewFileClientOpts
}

type NewFileClientOpts struct {
	// The base directory for all downloads.
	ClientBaseDir   string
	FilePathMaker   storage.FilePathMaker
	TorrentDirMaker storage.TorrentDirFilePathMaker
	PieceCompletion storage.PieceCompletion
}

// NewFileOpts creates a new ClientImplCloser that stores files using the OS native filesystem.
func NewFileOpts(opts NewFileClientOpts) storage.ClientImplCloser {
	if opts.TorrentDirMaker == nil {
		opts.TorrentDirMaker = defaultPathMaker
	}
	if opts.FilePathMaker == nil {
		opts.FilePathMaker = func(opts storage.FilePathMakerOpts) string {
			var parts []string
			if opts.Info.BestName() != metainfo.NoName {
				parts = append(parts, opts.Info.BestName())
			}
			return filepath.Join(append(parts, opts.File.BestPath()...)...)
		}
	}
	if opts.PieceCompletion == nil {
		opts.PieceCompletion = pieceCompletionForDir(opts.ClientBaseDir)
	}
	return fileClientImpl{opts}
}

func (me fileClientImpl) Close() error {
	return me.opts.PieceCompletion.Close()
}

func (fs fileClientImpl) OpenTorrent(
	ctx context.Context,
	info *metainfo.Info,
	infoHash metainfo.Hash,
) (_ storage.TorrentImpl, err error) {
	dir := fs.opts.TorrentDirMaker(fs.opts.ClientBaseDir, info, infoHash)
	logger := log.ContextLogger(ctx).Slogger()
	logger.DebugContext(ctx, "opened file torrent storage", slog.String("dir", dir))
	upvertedFiles := info.UpvertedFiles()
	files := make([]file, 0, len(upvertedFiles))
	for i, fileInfo := range upvertedFiles {
		filePath := filepath.Join(dir, fs.opts.FilePathMaker(storage.FilePathMakerOpts{
			Info: info,
			File: &fileInfo,
		}))
		if !isSubFilepath(dir, filePath) {
			err = fmt.Errorf("file %v: path %q is not sub path of %q", i, filePath, dir)
			return
		}
		f := file{
			path:   filePath,
			length: fileInfo.Length,
		}
		if f.length == 0 {
			err = CreateNativeZeroLengthFile(f.path)
			if err != nil {
				err = fmt.Errorf("creating zero length file: %w", err)
				return
			}
		}
		files = append(files, f)
	}
	t := &fileTorrentImpl{
		files,
		segments.NewIndexFromSegments(common.TorrentOffsetFileSegments(info)),
		infoHash,
		fs.opts.PieceCompletion,
	}
	return storage.TorrentImpl{
		Piece: t.Piece,
		Close: t.Close,
		Flush: t.Flush,
	}, nil
}

type file struct {
	// The safe, OS-local file path.
	path   string
	length int64
}

type fileTorrentImpl struct {
	files          []file
	segmentLocater segments.Index
	infoHash       metainfo.Hash
	completion     storage.PieceCompletion
}

func (fts *fileTorrentImpl) Piece(p metainfo.Piece) storage.PieceImpl {
	// Create a view onto the file-based torrent storage.
	_io := fileTorrentImplIO{fts}
	// Return the appropriate segments of this.
	return &filePieceImpl{
		fts,
		p,
		missinggo.NewSectionWriter(_io, p.Offset(), p.Length()),
		io.NewSectionReader(_io, p.Offset(), p.Length()),
	}
}

func (fs *fileTorrentImpl) Close() error {
	return nil
}

func fsync(filePath string) (err error) {
	_ = os.MkdirAll(filepath.Dir(filePath), 0o777)
	var f *os.File
	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = f.Sync(); err != nil {
		return err
	}
	return f.Close()
}

func (fts *fileTorrentImpl) Flush() error {
	for _, f := range fts.files {
		if err := fsync(f.path); err != nil {
			return err
		}
	}
	return nil
}

// A helper to create zero-length files which won't appear for file-orientated storage since no
// writes will ever occur to them (no torrent data is associated with a zero-length file). The
// caller should make sure the file name provided is safe/sanitized.
func CreateNativeZeroLengthFile(name string) error {
	os.MkdirAll(filepath.Dir(name), 0o777)
	var f io.Closer
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	return f.Close()
}

// Exposes file-based storage of a torrent, as one big ReadWriterAt.
type fileTorrentImplIO struct {
	fts *fileTorrentImpl
}

// Returns EOF on short or missing file.
func (fst *fileTorrentImplIO) readFileAt(file file, b []byte, off int64) (n int, err error) {
	f, err := os.Open(file.path)
	if os.IsNotExist(err) {
		// File missing is treated the same as a short file.
		err = io.EOF
		return
	}
	if err != nil {
		return
	}
	defer f.Close()
	// Limit the read to within the expected bounds of this file.
	if int64(len(b)) > file.length-off {
		b = b[:file.length-off]
	}
	for off < file.length && len(b) != 0 {
		n1, err1 := f.ReadAt(b, off)
		b = b[n1:]
		n += n1
		off += int64(n1)
		if n1 == 0 {
			err = err1
			break
		}
	}
	return
}

// Only returns EOF at the end of the torrent. Premature EOF is ErrUnexpectedEOF.
func (fst fileTorrentImplIO) ReadAt(b []byte, off int64) (n int, err error) {
	fst.fts.segmentLocater.Locate(segments.Extent{Start: off, Length: int64(len(b))}, func(i int, e segments.Extent) bool {
		n1, err1 := fst.readFileAt(fst.fts.files[i], b[:e.Length], e.Start)
		n += n1
		b = b[n1:]
		err = err1
		return err == nil // && int64(n1) == e.Length
	})
	if len(b) != 0 && err == nil {
		err = io.EOF
	}
	return
}

func (fst fileTorrentImplIO) WriteAt(p []byte, off int64) (n int, err error) {
	// log.Printf("write at %v: %v bytes", off, len(p))
	fst.fts.segmentLocater.Locate(segments.Extent{Start: off, Length: int64(len(p))}, func(i int, e segments.Extent) bool {
		name := fst.fts.files[i].path
		os.MkdirAll(filepath.Dir(name), 0o777)

		var f *os.File
		f, err = os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			return false
		}
		// defer os.Remove(f.Name())

		var n1 int
		n1, err = f.WriteAt(p[:e.Length], e.Start)
		// log.Println(fmt.Sprintf("%v %v wrote %v: %v", i, e, n1, err))
		closeErr := f.Close()
		n += n1
		p = p[n1:]
		if err == nil {
			err = closeErr
		}
		if err == nil && int64(n1) != e.Length {
			err = io.ErrShortWrite
		}
		return err == nil
	})
	return
}

// defaultPathMaker just returns the storage client's base directory.
func defaultPathMaker(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
	return baseDir
}

func isSubFilepath(base, sub string) bool {
	rel, err := filepath.Rel(base, sub)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func pieceCompletionForDir(dir string) (ret storage.PieceCompletion) {
	// This should be happening before sqlite attempts to open a database in the intended directory.
	os.MkdirAll(dir, 0o700)
	ret, err := storage.NewDefaultPieceCompletionForDir(dir)
	if err != nil {
		// This kinda sux using the global logger. This code is ancient.
		log.Levelf(log.Warning, "couldn't open piece completion db in %q: %s", dir, err)
		ret = storage.NewMapPieceCompletion()
	}
	return
}

type filePieceImpl struct {
	*fileTorrentImpl
	p metainfo.Piece
	io.WriterAt
	io.ReaderAt
}

func (me *filePieceImpl) pieceKey() metainfo.PieceKey {
	return metainfo.PieceKey{
		InfoHash: me.infoHash,
		Index:    me.p.Index(),
	}
}

func (fs *filePieceImpl) Completion() storage.Completion {
	c, err := fs.completion.Get(fs.pieceKey())
	if err != nil {
		log.Printf("error getting piece completion: %s", err)
		c.Ok = false
		return c
	}

	verified := true
	if c.Complete {
		// If it's allegedly complete, check that its constituent files have the necessary length.
		if !fs.segmentLocater.Locate(segments.Extent{
			Start:  fs.p.Offset(),
			Length: fs.p.Length(),
		}, func(i int, extent segments.Extent) bool {
			file := fs.files[i]
			s, err := os.Stat(file.path)
			if err != nil || s.Size() < extent.Start+extent.Length {
				verified = false
				return false
			}
			return true
		}) {
			panic("files do not cover piece extent")
		}
	}

	if !verified {
		// The completion was wrong, fix it.
		c.Complete = false
		fs.completion.Set(fs.pieceKey(), false)
	}

	return c
}

func (fs *filePieceImpl) MarkComplete() error {
	return fs.completion.Set(fs.pieceKey(), true)
}

func (fs *filePieceImpl) MarkNotComplete() error {
	return fs.completion.Set(fs.pieceKey(), false)
}
