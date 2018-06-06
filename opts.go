package txfile

import "fmt"

// Options provides common file options used when opening or creating a file.
type Options struct {
	// Additional flags.
	Flags Flag

	// MaxSize sets the maximum file size in bytes. This should be a multiple of PageSize.
	// If it's not a multiple of PageSize, the actual files maximum size is rounded downwards
	// to the next multiple of PageSize.
	// A value of 0 indicates the file can grow without limits.
	MaxSize uint64

	// PageSize sets the files page size on file creation. PageSize is ignored if
	// the file already exists.
	// If PageSize is not configured, the OSes main memory page size is selected.
	PageSize uint32

	// InitMetaArea configures the minimum amount of page in the meta area.
	// The amount of pages is only allocated when the file is generated.
	// The meta area grows by double the current meta area. To reduce the
	// total amount of pages moved to meta area on grow, it is recommended that
	// the value of InitMetaArea is a power of 2.
	InitMetaArea uint32

	// Prealloc disk space if MaxSize is set.
	Prealloc bool

	// Open file in readonly mode.
	Readonly bool
}

// Flag configures file opening behavior.
type Flag uint64

const (
	// FlagUnboundMaxSize configures the file max size to be unbound. This sets
	// MaxSize to 0. If MaxSize and Prealloc is set, up to MaxSize bytes are
	// preallocated on disk (truncate).
	FlagUnboundMaxSize Flag = 1 << iota

	// FlagUpdMaxSize updates the file max size setting. If not set, the max size
	// setting is read from the file to be opened.
	// The file will grow if MaxSize is larger then the current max size setting.
	// If MaxSize is less then the file's max size value, the file is tried to
	// shrink dynamically whenever pages are freed. Freed pages are returned via
	// `Truncate`.
	FlagUpdMaxSize
)

// Validate checks if all fields in Options are consistent with the File implementation.
func (o *Options) Validate() error {
	if o.Flags.Check(FlagUpdMaxSize) {
		if o.Readonly {
			return errReadOnlyUpdateSize
		}

		if !o.Flags.Check(FlagUnboundMaxSize) && o.MaxSize > 0 && o.MaxSize < minRequiredFileSize {
			return fmt.Errorf("max size must be at least %v bytes ", minRequiredFileSize)
		}
	}

	if metaSz := o.InitMetaArea; metaSz > 0 && o.MaxSize > 0 && o.PageSize > 0 {
		const headerPages = 2
		totalPages := o.MaxSize / uint64(o.PageSize)
		avail := totalPages - headerPages
		if uint64(metaSz) >= avail {
			return fmt.Errorf("meta area of %v pages exceeds the available pages %v",
				metaSz, avail)
		}
	}

	return nil
}

func (f Flag) Check(check Flag) bool {
	return (f & check) == check
}