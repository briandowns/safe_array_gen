#ifdef __cplusplus
extern "C" {
#endif

#ifndef __UTEST_H
#define __UTEST_H

#include <stdbool.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define UTEST_ASSERT_EQUAL(actual, expected) \
    do { \
        if (actual != expected) { \
            return false; \
        } \
    } while (0);

#define UTEST_ASSERT_EQUAL_SIZE_T  (actual, expected) \
    UTEST_ASSERT_EQUAL(actual,expected)
#define UTEST_ASSERT_EQUAL_UINT8_T (actual, expected) \
    UTEST_ASSERT_EQUAL(actual,expected)
// #define UTEST_ASSERT_EQUAL_UINT16_T(actual, expected) \
//     UTEST_ASSERT_EQUAL(actual,expected)
// #define UTEST_ASSERT_EQUAL_UINT32_T(actual, expected) \
//     UTEST_ASSERT_EQUAL(actual,expected)
// #define UTEST_ASSERT_EQUAL_UINT64_T(actual, expected) \
//     UTEST_ASSERT_EQUAL(actual,expected)
#define UTEST_ASSERT_EQUAL_INT     (actual, expected) \
    UTEST_ASSERT_EQUAL(actual,expected)

#define UTEST_ASSERT_EQUAL_STRING(actual, expected) \
    do { \
        if (strcmp(actual, expected) != 0) { \
            return false; \
        } \
    } while (0)

typedef enum {
    PASS,
    FAIL
} utest_status_t;

typedef bool (*utest_func_t)();

typedef struct {
    char name[256];
    clock_t start;
    clock_t end;
    utest_status_t status;
    utest_func_t func;
    char filename[];
} utest_t;

typedef struct {
    clock_t start;
    clock_t end;
    size_t count;
    size_t passed;
    size_t failed;
    utest_t *tests[512];
} utest_suite_t;

#define GREEN "\x1B[32m"
#define RED   "\x1B[31m"
#define RESET "\033[0m"

// #define UTEST_INIT \
//     test_suite.start = clock(); \
//     printf("\nRunning tests: %s\n", __FILE__);

void
utest_init();

void
utest_complete();

bool
utest_run(const char *name, utest_func_t func);

// #define RUN_TEST(name, func) { \
//     test_suite.count++; \
//     printf("    %-20s", name); \
//     clock_t test_start = clock(); \
//     func(); \
//     clock_t test_end = clock(); \
//     double time_spent = (double)(test_end - test_start) / CLOCKS_PER_SEC; \
//     test_suite.passed++; \
//     printf(GREEN "+"RESET"   %-2.3f/ms\n", (time_spent*1000)); \
// }

#define UTEST_COMPLETE \
    test_suite.end = clock(); \
    double ts = (double)(test_suite.end - test_suite.start) / CLOCKS_PER_SEC; \
    printf("\nTotal: %-4lu Passed: %-4lu Failed: %-4lu in  %-2.3f/ms\n", \
         test_suite.count, test_suite.passed, test_suite.failed, (ts*1000));

#endif /* end __UTEST_H */
#ifdef __cplusplus
}
#endif
