#ifdef __cplusplus
extern "C" {
#endif

#ifndef __CC_H
#define __CC_H

#include <stdbool.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define CC_ASSERT(condition) return condition
#define CC_ASSERT_EQUAL(actual, expected) \
    do { \
        if (actual != expected) { \
            return false; \
        } \
    } while (0)
#define CC_ASSERT_EQUAL_CHAR(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_STRING(actual, expected) \
    do { \
        if (strcmp(actual, expected) != 0) { \
            return false; \
        } \
    } while (0)
#define CC_ASSERT_EQUAL_SIZE_T(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_UINT8(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_UINT16(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_UINT32(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_UINT64(actual, expected) CC_ASSERT_EQUAL(actual, expected)
#define CC_ASSERT_EQUAL_INT(actual, expected) CC_ASSERT_EQUAL(actual, expected)

typedef bool (*utest_func_t)();

/**
 * Initializes the library and needed memory.
 */
void
utest_init();

/**
 * Cleans up used resources and prints results.
 */
void
utest_complete();

/**
 * Run the given test. This function can be called but needs to have the
 * additional arguments filled in to work. It's advised for users to use
 * the macro defined below.
 */
bool
utest_run(const char *name, utest_func_t func, const char *filename);

/**
 * Run the given test. This is the primary entry point.
 */
#define CC_RUN(name, func) { utest_run(name, func, __FILE__); }

#endif /* end __CC_H */
#ifdef __cplusplus
}
#endif
