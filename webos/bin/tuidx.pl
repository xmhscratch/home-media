#!/usr/bin/perl

use strict;
use warnings;
# use re 'debug';
use IPC::Open3;
use Symbol 'gensym';

use constant ASCII_RS => '\x{1E}';
use constant RGXP_MESSAGE_PAYLOAD => '^((\d+(?=\x{1E}))((?=\x{1E})..[^\x{1E}\n]*|)((?=\x{1E}).*)|.*)$';

sub trimRS {
    my $s = $_[0];
    $s =~ s/[@{[ASCII_RS]}]+//g;
    return $s;
}

my $msg = <>;
if ($msg =~ m/@{[RGXP_MESSAGE_PAYLOAD]}/gs) {
    my $cmdStr = trimRS($4);
    # open(my $fh, '-|', ) or die $!;

    my $errfh = gensym;  # separate handle for STDERR
    my $pid = open3(my $infh, my $outfh, $errfh, $cmdStr);
    # (Writes to $infh, reads from $outfh and $errfh)
    close($infh);  # if no input to send

    my $stdout = do { local $/; <$outfh> };
    my $stderr = do { local $/; <$errfh> };

    waitpid($pid, 0);
    my $exit_status = $? >> 8;
    # print $?;
    # print $exit_status;
    print $stdout;
    print $stderr;
    # print $outfh;

    # exec '/usr/bin/go run', "/home/web/repos/home-media/cmd/tui/main.go -m=text $stdout";
    # open(@message, '>' , '/home/web/repos/home-media/test.txt') or die $!;
    # close(@message);
}
