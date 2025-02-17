#include <stdbool.h>
#include <stdlib.h>
#include <time.h>

#include "utest.h"

#define GREEN "\x1B[32m"
#define RED   "\x1B[31m"
#define RESET "\033[0m"

static size_t count = 0;
static size_t passed = 0;
static size_t failed = 0;
static clock_t start = 0;
static clock_t end = 0;

void
utest_init()
{
    start = clock();
    printf("Running test suite...\n");
}

bool
utest_run(const char *name, utest_func_t func, const char *func_name) {
    count++;

    clock_t test_start = clock();
    bool ret = func();
    clock_t test_end = clock();

    double time_spent = (double)(test_end - test_start) / CLOCKS_PER_SEC;
    if (ret == false) {
        failed++;
        printf("    %-28s%-2s:%-21d" RED "%-8s" RESET "   %-2.3f/ms\n",
            name, func_name, 0, "failed", (time_spent*1000));
        return false;
    }

    passed++;
    printf("    %-28s%-28s" GREEN "%-8s" RESET "   %-2.3f/ms\n",
         name, "", "passed", (time_spent*1000));

    return true;
}

void
utest_complete()
{
    end = clock();
    double ts = (double)(end - start) / CLOCKS_PER_SEC; \
    printf("\nTotal: %-4lu Passed: %-4lu Failed: %-4lu in  %-2.3f/ms\n", \
         count, passed, failed, (ts*1000));
}
