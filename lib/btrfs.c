#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mount.h>
#include "mosaic.h"
#include "tessera.h"

static int open_btrfs(struct mosaic *m, int flags)
{
	m->default_fs = strdup("btrfs");
	return 0;
}

static int new_btrfs_subvol(struct mosaic *m, char *name,
		unsigned long size_in_blocks, int make_flags)
{
	char aux[1024];

	/*
	 * FIXME: locate tesserae subvolumes in subdirectories
	 * FIXME: what to do with "size_in_blocks"? Tree quota?
	 */
	sprintf(aux, "btrfs subvolume create %s/%s", m->m_loc, name);
	if (system(aux))
		return -1;

	return 0;
}

static int open_btrfs_subvol(struct mosaic *m, struct tessera *t,
		int open_flags)
{
	return 0; /* FIXME: check it exists */
}

static int clone_btrfs_subvol(struct mosaic *m, struct tessera *from,
		char *name, int clone_flags)
{
	char aux[1024];

	/*
	 * FIXME: locate tesserae subvolumes in subdirectories
	 */
	sprintf(aux, "btrfs subvolume create %s/%s %s/%s",
			m->m_loc, from->t_name, m->m_loc, name);
	if (system(aux))
		return -1;

	return 0;
}

static int drop_btrfs_subvol(struct mosaic *m, struct tessera *t,
		int drop_flags)
{
	char aux[1024];

	/*
	 * FIXME: locate tesserae subvolumes in subdirectories
	 */
	sprintf(aux, "btrfs subvolume delete %s/%s", m->m_loc, t->t_name);
	if (system(aux))
		return -1;

	return 0;
}

static int mount_btrfs_subvol(struct mosaic *m, struct tessera *t,
		char *path, int mount_flags)
{
	char aux[1024];

	sprintf(aux, "%s/%s", m->m_loc, t->t_name);
	return mount(aux, path, NULL, MS_BIND | mount_flags, NULL);
}

static int resize_btrfs_subvol(struct mosaic *m, struct tessera *t,
		unsigned long size_in_blocks, int resize_flags)
{
	return -1;
}

const struct mosaic_ops mosaic_btrfs = {
	.open = open_btrfs,
	.mount = bind_mosaic_loc,

	.new_tessera = new_btrfs_subvol,
	.open_tessera = open_btrfs_subvol,
	.clone_tessera = clone_btrfs_subvol,
	.drop_tessera = drop_btrfs_subvol,
	.mount_tessera = mount_btrfs_subvol,
	.resize_tessera = resize_btrfs_subvol,
};
