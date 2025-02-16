#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <time.h>

#include "int_slice.h"

void
print_item(int item, void *user_data)
{
    printf("%d\n", item);
}

bool
compare_func(const int x, const int y, void *user_data)
{
    return x == y;
}

void
test_slice_append()
{
    int_slice_t *s1 = int_slice_new(8);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    int_slice_append(s1, 88);
    int_slice_append(s1, 70);
    int_slice_append(s1, 41);
    int_slice_append(s1, 15);
    int_slice_append(s1, 1);

    assert(s1->len == 8);
    assert(s1->cap == 8);

    int_slice_free(s1);
}

void
test_slice_get()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);

    assert(int_slice_get(s1, 0) == 9);
    assert(int_slice_get(s1, 1) == 100);
    assert(int_slice_get(s1, 2) == 7);

    int_slice_free(s1);
}

void
test_slice_contains()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    
    assert(int_slice_contains(s1, 9));
    assert(int_slice_contains(s1, 100));
    assert(int_slice_contains(s1, 7));

    int_slice_free(s1);
}

void
test_slice_delete()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    
    int_slice_delete(s1, 2);

    assert(s1->len == 2);

    int_slice_free(s1);
}

void
test_slice_replace_by_index()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_replace_by_idx(s1, 0, 42);

    assert(int_slice_get(s1, 0) == 42);

    int_slice_free(s1);
}

void
test_slice_count()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);

    size_t count = int_slice_count(s1, 7);
    assert(count == 3);

    int_slice_free(s1);
}

void
test_slice_sort()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 100);
    int_slice_append(s1, 1000);
    int_slice_append(s1, 9);

    int_slice_sort(s1);
    assert(int_slice_get(s1, 0) == 9);
    assert(int_slice_get(s1, 1) == 100);
    assert(int_slice_get(s1, 2) == 1000);

    int_slice_free(s1);
}

void
test_slice_replace_by_val()
{
    int_slice_t *s1 = int_slice_new(3);
    s1->compare = compare_func;

    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);

    int_slice_replace_by_val(s1, 100, 700, 3);

    assert(int_slice_count(s1, 700) == 3);
    assert(int_slice_count(s1, 100) == 1);

    int_slice_free(s1);
}

void
test_slice_repeat()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_repeat(s1, 88, 10);

    assert(s1->len == 10);
    assert(int_slice_get(s1, 0) == 88);
    assert(int_slice_get(s1, s1->len-1) == 88);

    int_slice_free(s1);
}

void
test_slice_grow()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_grow(s1, 2);
    assert(s1->cap == 5);

    int_slice_free(s1);
}

void
test_slice_concat()
{
    int_slice_t *s1 = int_slice_new(3);
    int_slice_t *s2 = int_slice_new(3);

    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);

    assert(s1->len == 3);
    assert(s1->cap == 3);

    int_slice_append(s2, 9);
    int_slice_append(s2, 9);
    int_slice_append(s2, 9);

    assert(s2->len == 3);
    assert(s2->cap == 3);

    size_t s1_len = int_slice_concat(s1, s2);
    assert(s1_len == 6);

    int_slice_free(s1);
    int_slice_free(s2);
}

typedef void (*test_func_t)();

typedef struct {
    clock_t start;
    clock_t end;
    test_func_t *tests;
} test_suite_t;

static test_suite_t test_suite = {0};

#define TESTS_INIT \
    test_suite.start = clock(); \
    printf("\nRunning tests: %s\n", __FILE__);

#define RUN_TEST(name, func) { \
    printf("    %-20s", name); \
    clock_t test_start = clock(); \
    func(); \
    clock_t test_end = clock(); \
    double time_spent = (double)(test_end - test_start) / CLOCKS_PER_SEC; \
    printf("successful   %-2.3f/ms\n", (time_spent*1000)); \
}

#define TESTS_COMPLETE \
    test_suite.end = clock(); \
    double ts = (double)(test_suite.end - test_suite.start) / CLOCKS_PER_SEC; \
    printf("\nComplete: %-2.3f/ms\n", (ts*1000));

int
main(int argc, char **argv)
{
    TESTS_INIT;

    RUN_TEST("append", test_slice_append);
    RUN_TEST("get", test_slice_get);
    RUN_TEST("contains", test_slice_contains);
    RUN_TEST("delete", test_slice_delete);
    RUN_TEST("replace by index", test_slice_replace_by_index);
    RUN_TEST("count", test_slice_count);
    RUN_TEST("sort", test_slice_sort);
    RUN_TEST("replace by value", test_slice_replace_by_val);
    RUN_TEST("repeat", test_slice_repeat);
    RUN_TEST("grow", test_slice_grow);
    RUN_TEST("concat", test_slice_concat);

    TESTS_COMPLETE;

    return 0;
}
