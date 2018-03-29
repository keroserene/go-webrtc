// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef BASE_TASK_SCHEDULER_TEST_UTILS_H_
#define BASE_TASK_SCHEDULER_TEST_UTILS_H_

#include <memory>

#include "base/memory/ref_counted.h"
#include "base/task_runner.h"

#include "base/task_scheduler/sequence.h"

namespace base {
namespace internal {

class SchedulerWorkerPool;
struct Task;

namespace test {

// An enumeration of possible task scheduler TaskRunner types. Used to
// parametrize relevant task_scheduler tests.
enum class ExecutionMode { PARALLEL, SEQUENCED, SINGLE_THREADED };

// Creates a Sequence and pushes |task| to it. Returns that sequence.
scoped_refptr<Sequence> CreateSequenceWithTask(std::unique_ptr<Task> task);

// Creates a TaskRunner that posts tasks to |worker_pool| with the
// |execution_mode| execution mode and the WithBaseSyncPrimitives() trait.
// Caveat: this does not support ExecutionMode::SINGLE_THREADED.
scoped_refptr<TaskRunner> CreateTaskRunnerWithExecutionMode(
    SchedulerWorkerPool* worker_pool,
    test::ExecutionMode execution_mode);

}  // namespace test
}  // namespace internal
}  // namespace base

#endif  // BASE_TASK_SCHEDULER_TEST_UTILS_H_
