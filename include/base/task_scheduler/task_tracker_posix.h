// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef BASE_TASK_SCHEDULER_TASK_TRACKER_POSIX_H_
#define BASE_TASK_SCHEDULER_TASK_TRACKER_POSIX_H_

#include <memory>

#include "base/base_export.h"
#include "base/logging.h"
#include "base/macros.h"
#include "base/task_scheduler/task_tracker.h"
#include "base/threading/platform_thread.h"

namespace base {

class MessageLoopForIO;

namespace internal {

struct Task;

// A TaskTracker that instantiates a FileDescriptorWatcher in the scope in which
// a task runs. Used on all POSIX platforms except NaCl SFI.
// set_watch_file_descriptor_message_loop() must be called before the
// TaskTracker can run tasks.
class BASE_EXPORT TaskTrackerPosix : public TaskTracker {
 public:
  TaskTrackerPosix();
  ~TaskTrackerPosix() override;

  // Sets the MessageLoopForIO with which to setup FileDescriptorWatcher in the
  // scope in which tasks run. Must be called before starting to run tasks.
  // External synchronization is required between a call to this and a call to
  // RunTask().
  void set_watch_file_descriptor_message_loop(
      MessageLoopForIO* watch_file_descriptor_message_loop) {
    watch_file_descriptor_message_loop_ = watch_file_descriptor_message_loop;
  }

#if DCHECK_IS_ON()
  // TODO(robliao): http://crbug.com/698140. This addresses service thread tasks
  // that could run after the task scheduler has shut down. Anything from the
  // service thread is exempted from the task scheduler shutdown DCHECKs.
  void set_service_thread_handle(
      const PlatformThreadHandle& service_thread_handle) {
    DCHECK(!service_thread_handle.is_null());
    service_thread_handle_ = service_thread_handle;
  }
#endif

 protected:
  // TaskTracker:
  void RunOrSkipTask(std::unique_ptr<Task> task,
                     Sequence* sequence,
                     bool can_run_task) override;

 private:
#if DCHECK_IS_ON()
  bool IsPostingBlockShutdownTaskAfterShutdownAllowed() override;
#endif

  MessageLoopForIO* watch_file_descriptor_message_loop_ = nullptr;

#if DCHECK_IS_ON()
  PlatformThreadHandle service_thread_handle_;
#endif

  DISALLOW_COPY_AND_ASSIGN(TaskTrackerPosix);
};

}  // namespace internal
}  // namespace base

#endif  // BASE_TASK_SCHEDULER_TASK_TRACKER_POSIX_H_
