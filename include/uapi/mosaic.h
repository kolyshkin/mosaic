#ifndef __MOSAIC_UAPI_H__
#define __MOSAIC_UAPI_H__

#pragma GCC visibility push(default)

/*
 * Mosaic management
 */

typedef struct mosaic *mosaic_t;
mosaic_t mosaic_open(const char *name, int open_flags);
void mosaic_close(mosaic_t m);
int mosaic_get_name(mosaic_t m, char *name_buf, int buf_len);
int mosaic_mount(mosaic_t m, const char *path, int mount_flags);

/*
 * Volume management
 */

typedef struct volume *volume_t;
volume_t mosaic_open_vol(mosaic_t m, const char *name, int open_flags);
void mosaic_close_vol(volume_t);

/** Checks if a volume with a given @name exists in mosaic @m.
 * Argument @flags is currently unused and should be 0.
 *
 * Returns:
 *  1: exists
 *  0: does not exist
 * -1: some error
 */
int mosaic_have_vol(mosaic_t m, const char *name, int flags);

/*
 * Make a new volume. The first one makes raw block device, the
 * second one also puts filesystem on it.
 */
int mosaic_make_vol(mosaic_t m, const char *name, unsigned long size_in_blocks, int make_flags);
int mosaic_make_vol_fs(mosaic_t m, const char *name, unsigned long size_in_blocks, int make_flags);

/*
 * Create a clone of existing volume with COW (when possible)
 */
int mosaic_clone_vol(volume_t from, const char *name, int clone_flags);

int mosaic_drop_vol(volume_t t, int drop_flags);
int mosaic_resize_vol(volume_t t, unsigned long new_size_in_blocks, int resize_flags);

/*
 * Mounting and umounting of volume
 */
int mosaic_mount_vol(volume_t t, const char *path, int mount_flags);
int mosaic_umount_vol(volume_t t, const char *path, int umount_flags);

/*
 * Getting path to volume's block device (when possible)
 * and putting it back.
 *
 * Return value from the first one is the lenght of the name
 * of the device (even if it doesn't fit the buffer len).
 */
int mosaic_get_vol_bdev(volume_t t, char *dev, int len, int flags);
int mosaic_put_vol_bdev(volume_t t);

int mosaic_get_vol_size(volume_t t, unsigned long *size_in_blocks);

/*
 * (Live) migration
 */

int mosaic_migrate_vol_send_start(volume_t t, int fd_to, int flags);
int mosaic_migrate_vol_send_more(volume_t t);
int mosaic_migrate_vol_recv_start(volume_t t, int fd_from, int flags);
int mosaic_migrate_vol_stop(volume_t t);

/* Misc */

/* Can do mosaic_clone_vol */
#define MOSAIC_FEATURE_CLONE		(1ULL << 0)
/* The size_in_blocks works */
#define MOSAIC_FEATURE_DISK_SIZE_MGMT	(1ULL << 1)
/* Support block device opening */
#define MOSAIC_FEATURE_BDEV		(1ULL << 2)
/* Supports iterative copying */
#define MOSAIC_FEATURE_MIGRATE		(1ULL << 3)

int mosaic_get_features(mosaic_t mos, unsigned long long *feat);

typedef void (*mosaic_log_fn)(int lvl, const char *f, ...)
	__attribute__ ((format(printf, 2, 3)));
void mosaic_set_log_fn(mosaic_log_fn lfn);

enum log_level {
	LOG_ERR = 1,
	LOG_WRN,
	LOG_INF,
	LOG_DBG,
};

void mosaic_set_log_lvl(enum log_level l);

#pragma GCC visibility pop
#endif
