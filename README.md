dwmstatus
============

Fork off oniichaNj/go-dwmstatus
Current Changes from parent:
 * removed battery check, since I'm using a desktop, might add this back in with graceful failure if I ever care to distribute
 * removed alsa-bound getvol.h C functions
 * added pulseaudio volume grabber, taken from github.com/vially/volumectl/pulseaudio
 * -again, ideally this would gracefully detect what kind of sound you're using, but currently this is personal use only
 * likely to be more personal touches
