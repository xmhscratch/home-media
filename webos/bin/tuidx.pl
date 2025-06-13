#!/usr/bin/perl

use strict;
use warnings;
# use re 'debug';

use constant ASCII_RS => '\x{1E}';
use constant RGXP_MESSAGE_PAYLOAD => '^((\d+(?=\x{1E}))((?=\x{1E})..[^\x{1E}\n]*|)((?=\x{1E}).*)|.*)$';

sub trimRS {
    my $s = $_[0];
    $s =~ s/[@{[ASCII_RS]}]+//g;
    return $s;
}

# open my $outFile, '>', '/home/web/repos/home-media/test.txt';
# print $outFile `$cmdStr`;
# close $outFile;

my $msg = <>;
if ($msg =~ m/@{[RGXP_MESSAGE_PAYLOAD]}/gs) {
    my $cmdStr = trimRS($4);
    exec '/bin/echo', `$cmdStr`;
}
