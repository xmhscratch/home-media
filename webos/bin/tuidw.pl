#!/usr/bin/perl

use strict;
use warnings;
use IPC::Open3;
use Symbol 'gensym';
use IO::Socket::UNIX;

use constant ASCII_RS => '\x{1E}';
use constant UNIX_VW_SOCKET_PATH => '/run/tuidw.sock';

# apk del font-noto font-noto-common font-noto-emoji font-noto-extra

use IPC::Open3;
use Symbol qw(gensym);
use IO::Handle;
use IO::Select;

my $cmd = <STDIN>;
my $errfh = gensym;
my $pid = open3(my $infh, my $outfh, $errfh, $cmd);
close($infh);  # if no input to send

$outfh->autoflush(1);
$errfh->autoflush(1);

my $sel = IO::Select->new();
$sel->add($outfh, $errfh);

sub sendMessage {
    my $mode = $_[0];
    my $opts = $_[1];
    my $message = $_[2];

    my $sock = IO::Socket::UNIX->new(
        Type => SOCK_STREAM(),
        Peer => UNIX_VW_SOCKET_PATH,
    ) or die "Can't connect to @{[UNIX_VW_SOCKET_PATH]}: $!";

    $message =~ s/(['])/\\$1/g;

    print $sock "$mode@{[ASCII_RS]}$opts@{[ASCII_RS]}$message";
    print $message;

    $sock->close;
    return $? >> 8;
}

# Read and print stdout in real-time
while (my $line = <$outfh>) {
    next unless defined $line;
    sendMessage(6,"",$line);
}

# Optionally: read stderr after stdout ends
while (my $line = <$errfh>) {
    sendMessage(3,"",$line);
}

waitpid($pid, 0);

my $exit_code = $? >> 8;
if ($? == -1) {
    sendMessage(3,"","Failed to execute: $!\n");
} elsif ($? & 127) {
    sendMessage(3,"",sprintf("Died with signal %d\n", $? & 127));
} else {
    sendMessage(3,"","Process exited with code $exit_code\n");
}
