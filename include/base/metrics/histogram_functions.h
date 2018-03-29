// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef BASE_METRICS_HISTOGRAM_FUNCTIONS_H_
#define BASE_METRICS_HISTOGRAM_FUNCTIONS_H_

#include "base/metrics/histogram.h"
#include "base/metrics/histogram_base.h"
#include "base/time/time.h"

// Functions for recording metrics.
//
// For best practices on deciding when to emit to a histogram and what form
// the histogram should take, see
// https://chromium.googlesource.com/chromium/src.git/+/HEAD/tools/metrics/histograms/README.md

// Functions for recording UMA histograms. These can be used for cases
// when the histogram name is generated at runtime. The functionality is
// equivalent to macros defined in histogram_macros.h but allowing non-constant
// histogram names. These functions are slower compared to their macro
// equivalent because the histogram objects are not cached between calls.
// So, these shouldn't be used in performance critical code.
namespace base {

// For histograms with linear buckets.
// Used for capturing integer data with a linear bucketing scheme. This can be
// used when you want the exact value of some small numeric count, with a max of
// 100 or less. If you need to capture a range of greater than 100, we recommend
// the use of the COUNT histograms below.
// Sample usage:
//   base::UmaHistogramExactLinear("Histogram.Linear", some_value, 10);
BASE_EXPORT void UmaHistogramExactLinear(const std::string& name,
                                         int sample,
                                         int value_max);

// For adding a sample to an enumerated histogram.
// Sample usage:
//   // These values are persisted to logs. Entries should not be renumbered and
//   // numeric values should never be reused.
//   enum class MyEnum {
//     FIRST_VALUE = 0,
//     SECOND_VALUE = 1,
//     ...
//     FINAL_VALUE = N,
//     COUNT
//   };
//   base::UmaHistogramEnumeration("My.Enumeration",
//                                 MyEnum::SOME_VALUE, MyEnum::COUNT);
//
// Note: The value in |sample| must be strictly less than |enum_size|.
template <typename T>
void UmaHistogramEnumeration(const std::string& name, T sample, T enum_size) {
  static_assert(std::is_enum<T>::value,
                "Non enum passed to UmaHistogramEnumeration");
  DCHECK_LE(static_cast<uintmax_t>(enum_size), static_cast<uintmax_t>(INT_MAX));
  DCHECK_LE(static_cast<uintmax_t>(sample), static_cast<uintmax_t>(enum_size));
  return UmaHistogramExactLinear(name, static_cast<int>(sample),
                                 static_cast<int>(enum_size));
}

// For adding boolean sample to histogram.
// Sample usage:
//   base::UmaHistogramBoolean("My.Boolean", true)
BASE_EXPORT void UmaHistogramBoolean(const std::string& name, bool sample);

// For adding histogram with percent.
// Percents are integer between 1 and 100.
// Sample usage:
//   base::UmaHistogramPercentage("My.Percent", 69)
BASE_EXPORT void UmaHistogramPercentage(const std::string& name, int percent);

// For adding counts histogram.
// Sample usage:
//   base::UmaHistogramCustomCounts("My.Counts", some_value, 1, 600, 30)
BASE_EXPORT void UmaHistogramCustomCounts(const std::string& name,
                                          int sample,
                                          int min,
                                          int max,
                                          int buckets);

// Counts specialization for maximum counts 100, 1000, 10k, 100k, 1M and 10M.
BASE_EXPORT void UmaHistogramCounts100(const std::string& name, int sample);
BASE_EXPORT void UmaHistogramCounts1000(const std::string& name, int sample);
BASE_EXPORT void UmaHistogramCounts10000(const std::string& name, int sample);
BASE_EXPORT void UmaHistogramCounts100000(const std::string& name, int sample);
BASE_EXPORT void UmaHistogramCounts1M(const std::string& name, int sample);
BASE_EXPORT void UmaHistogramCounts10M(const std::string& name, int sample);

// For histograms storing times.
BASE_EXPORT void UmaHistogramCustomTimes(const std::string& name,
                                         TimeDelta sample,
                                         TimeDelta min,
                                         TimeDelta max,
                                         int buckets);
// For short timings from 1 ms up to 10 seconds (50 buckets).
BASE_EXPORT void UmaHistogramTimes(const std::string& name, TimeDelta sample);
// For medium timings up to 3 minutes (50 buckets).
BASE_EXPORT void UmaHistogramMediumTimes(const std::string& name,
                                         TimeDelta sample);
// For time intervals up to 1 hr (50 buckets).
BASE_EXPORT void UmaHistogramLongTimes(const std::string& name,
                                       TimeDelta sample);

// For recording memory related histograms.
// Used to measure common KB-granularity memory stats. Range is up to 500M.
BASE_EXPORT void UmaHistogramMemoryKB(const std::string& name, int sample);
// Used to measure common MB-granularity memory stats. Range is up to ~1G.
BASE_EXPORT void UmaHistogramMemoryMB(const std::string& name, int sample);
// Used to measure common MB-granularity memory stats. Range is up to ~64G.
BASE_EXPORT void UmaHistogramMemoryLargeMB(const std::string& name, int sample);

}  // namespace base

#endif  // BASE_METRICS_HISTOGRAM_FUNCTIONS_H_
