#ifndef __MOSAIC_MOCTL_H__
#define __MOSAIC_MOCTL_H__
struct mosaic;

int create_fsimg(struct mosaic *, int argc, char **argv);
int create_plain(struct mosaic *, int argc, char **argv);
int create_btrfs(struct mosaic *, int argc, char **argv);
#endif