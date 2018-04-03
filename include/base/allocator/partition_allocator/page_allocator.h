// Copyright (c) 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef BASE_ALLOCATOR_PARTITION_ALLOCATOR_PAGE_ALLOCATOR_H
#define BASE_ALLOCATOR_PARTITION_ALLOCATOR_PAGE_ALLOCATOR_H

#include <stdint.h>

#include <cstddef>

#include "base/base_export.h"
#include "base/compiler_specific.h"
#include "build/build_config.h"

namespace base {

#if defined(OS_WIN)
static const size_t kPageAllocationGranularityShift = 16;  // 64KB
#else
static const size_t kPageAllocationGranularityShift = 12;  // 4KB
#endif
static const size_t kPageAllocationGranularity =
    1 << kPageAllocationGranularityShift;
static const size_t kPageAllocationGranularityOffsetMask =
    kPageAllocationGranularity - 1;
static const size_t kPageAllocationGranularityBaseMask =
    ~kPageAllocationGranularityOffsetMask;

// All Blink-supported systems have 4096 sized system pages and can handle
// permissions and commit / decommit at this granularity.
static const size_t kSystemPageSize = 4096;
static const size_t kSystemPageOffsetMask = kSystemPageSize - 1;
static_assert((kSystemPageSize & (kSystemPageSize - 1)) == 0,
              "kSystemPageSize must be power of 2");
static const size_t kSystemPageBaseMask = ~kSystemPageOffsetMask;

enum PageAccessibilityConfiguration {
  PageInaccessible,
  PageReadWrite,
  PageReadExecute,
  PageReadWriteExecute,
};

// Allocate one or more pages.
//
// The requested |address| is just a hint; the actual address returned may
// differ. The returned address will be aligned at least to |align| bytes.
// |length| is in bytes, and must be a multiple of |kPageAllocationGranularity|.
// |align| is in bytes, and must be a power-of-two multiple of
// |kPageAllocationGranularity|.
//
// If |address| is null, then a suitable and randomized address will be chosen
// automatically.
//
// |page_accessibility| controls the permission of the allocated pages.
//
// This call will return null if the allocation cannot be satisfied.
BASE_EXPORT void* AllocPages(void* address,
                             size_t length,
                             size_t align,
                             PageAccessibilityConfiguration page_accessibility,
                             bool commit = true);

// Free one or more pages starting at |address| and continuing for |length|
// bytes.
//
// |address| and |length| must match a previous call to |AllocPages|. Therefore,
// |address| must be aligned to |kPageAllocationGranularity| bytes, and |length|
// must be a multiple of |kPageAllocationGranularity|.
BASE_EXPORT void FreePages(void* address, size_t length);

// Mark one or more system pages, starting at |address| with the given
// |page_accessibility|. |length| must be a multiple of |kSystemPageSize| bytes.
//
// Returns true if the permission change succeeded. In most cases you must
// |CHECK| the result.
BASE_EXPORT WARN_UNUSED_RESULT bool SetSystemPagesAccess(
    void* address,
    size_t length,
    PageAccessibilityConfiguration page_accessibility);

// Decommit one or more system pages starting at |address| and continuing for
// |length| bytes. |length| must be a multiple of |kSystemPageSize|.
//
// Decommitted means that physical resources (RAM or swap) backing the allocated
// virtual address range are released back to the system, but the address space
// is still allocated to the process (possibly using up page table entries or
// other accounting resources). Any access to a decommitted region of memory
// is an error and will generate a fault.
//
// This operation is not atomic on all platforms.
//
// Note: "Committed memory" is a Windows Memory Subsystem concept that ensures
// processes will not fault when touching a committed memory region. There is
// no analogue in the POSIX memory API where virtual memory pages are
// best-effort allocated resources on the first touch. To create a
// platform-agnostic abstraction, this API simulates the Windows "decommit"
// state by both discarding the region (allowing the OS to avoid swap
// operations) and changing the page protections so accesses fault.
//
// TODO(ajwong): This currently does not change page protections on POSIX
// systems due to a perf regression. Tracked at http://crbug.com/766882.
BASE_EXPORT void DecommitSystemPages(void* address, size_t length);

// Recommit one or more system pages, starting at |address| and continuing for
// |length| bytes with the given |page_accessibility|. |length| must be a
// multiple of |kSystemPageSize|.
//
// Decommitted system pages must be recommitted with their original permissions
// before they are used again.
//
// Returns true if the recommit change succeeded. In most cases you must |CHECK|
// the result.
BASE_EXPORT WARN_UNUSED_RESULT bool RecommitSystemPages(
    void* address,
    size_t length,
    PageAccessibilityConfiguration page_accessibility);

// Discard one or more system pages starting at |address| and continuing for
// |length| bytes. |length| must be a multiple of |kSystemPageSize|.
//
// Discarding is a hint to the system that the page is no longer required. The
// hint may:
//   - Do nothing.
//   - Discard the page immediately, freeing up physical pages.
//   - Discard the page at some time in the future in response to memory
//   pressure.
//
// Only committed pages should be discarded. Discarding a page does not decommit
// it, and it is valid to discard an already-discarded page. A read or write to
// a discarded page will not fault.
//
// Reading from a discarded page may return the original page content, or a page
// full of zeroes.
//
// Writing to a discarded page is the only guaranteed way to tell the system
// that the page is required again. Once written to, the content of the page is
// guaranteed stable once more. After being written to, the page content may be
// based on the original page content, or a page of zeroes.
BASE_EXPORT void DiscardSystemPages(void* address, size_t length);

ALWAYS_INLINE uintptr_t RoundUpToSystemPage(uintptr_t address) {
  return (address + kSystemPageOffsetMask) & kSystemPageBaseMask;
}

ALWAYS_INLINE uintptr_t RoundDownToSystemPage(uintptr_t address) {
  return address & kSystemPageBaseMask;
}

// Reserves (at least) |size| bytes of address space, aligned to
// |kPageAllocationGranularity|. This can be called early on to make it more
// likely that large allocations will succeed. Returns true if the reservation
// succeeded, false if the reservation failed or a reservation was already made.
BASE_EXPORT bool ReserveAddressSpace(size_t size);

// Releases any reserved address space. |AllocPages| calls this automatically on
// an allocation failure. External allocators may also call this on failure.
BASE_EXPORT void ReleaseReservation();

// Returns |errno| (POSIX) or the result of |GetLastError| (Windows) when |mmap|
// (POSIX) or |VirtualAlloc| (Windows) fails.
BASE_EXPORT uint32_t GetAllocPageErrorCode();

}  // namespace base

#endif  // BASE_ALLOCATOR_PARTITION_ALLOCATOR_PAGE_ALLOCATOR_H
