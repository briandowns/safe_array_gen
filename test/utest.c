#include <stdbool.h>
#include <stdlib.h>
#include <time.h>

#include "utest.h"

static utest_suite_t *test_suite = {0};

void
utest_init()
{
    test_suite = (utest_suite_t*)calloc(1, sizeof(utest_suite_t));

    test_suite->start = clock(); \
    printf("\nRunning tests: %s\n", __FILE__);
}

bool
utest_run(const char *name, utest_func_t func) {
    test_suite->count++;

    printf("    %-20s", name);

    clock_t test_start = clock();
    bool ret = func();
    clock_t test_end = clock();

    double time_spent = (double)(test_end - test_start) / CLOCKS_PER_SEC;
    if (ret) {
        test_suite->failed++;
        printf(RED"%-8s"RESET"   %-2.3f/ms\n", "failed", (time_spent*1000));
        return false;
    }

    test_suite->passed++;
    printf(GREEN "%-8s"RESET"   %-2.3f/ms\n", "passed", (time_spent*1000));

    return true;
}

void
utest_complete()
{
    test_suite->end = clock();
    double ts = (double)(test_suite->end - test_suite->start) / CLOCKS_PER_SEC; \
    printf("\nTotal: %-4lu Passed: %-4lu Failed: %-4lu in  %-2.3f/ms\n", \
         test_suite->count, test_suite->passed, test_suite->failed, (ts*1000));
    free(test_suite);
}
