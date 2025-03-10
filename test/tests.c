#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "int_slice.h"
#include "crosscheck.h"

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

int
qsort_compare(const void *x, const void *y) {
	return (*(int*)x - *(int*)y);
}

int_slice_t *s1;

void
cc_setup()
{
    s1 = int_slice_new(3);
}

void
cc_tear_down()
{
    int_slice_free(s1);
}

cc_result_t
test_slice_append()
{
    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);

    CC_ASSERT_UINT64_EQUAL(s1->len, (uint64_t)3);
    CC_ASSERT_UINT64_EQUAL(s1->cap, (uint64_t)3);

    CC_SUCCESS;
}

cc_result_t
test_slice_get()
{
    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);

    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 0), 9);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 1), 100);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 2), 7);

    CC_SUCCESS;
}

cc_result_t
test_slice_contains()
{
    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    
    CC_ASSERT_TRUE(int_slice_contains(s1, 9));
    CC_ASSERT_TRUE(int_slice_contains(s1, 100));
    CC_ASSERT_TRUE(int_slice_contains(s1, 7));

    CC_SUCCESS;
}

cc_result_t
test_slice_delete()
{
    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    int_slice_delete(s1, 2);

    CC_ASSERT_UINT64_EQUAL(s1->len, 2);

    CC_SUCCESS;
}

cc_result_t
test_slice_replace_by_index()
{
    int_slice_append(s1, 9);
    int_slice_replace_by_idx(s1, 0, 42);

    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 0), 42);

    CC_SUCCESS;
}

cc_result_t
test_slice_count()
{
    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    uint64_t count = int_slice_count(s1, 7, compare_func);

    CC_ASSERT_UINT64_EQUAL(count, 3);

    CC_SUCCESS;
}

cc_result_t
test_slice_sort()
{
    int_slice_append(s1, 100);
    int_slice_append(s1, 1000);
    int_slice_append(s1, 9);
    int_slice_sort(s1, qsort_compare);

    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 0), 9);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 1), 100);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 2), 1000);

    CC_SUCCESS;
}

cc_result_t
test_slice_replace_by_val()
{
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_replace_by_val(s1, 100, 700, 3, compare_func);

    CC_ASSERT_UINT64_EQUAL(int_slice_count(s1, 700, compare_func), 3);
    CC_ASSERT_UINT64_EQUAL(int_slice_count(s1, 100, compare_func), 1);

    CC_SUCCESS;
}

cc_result_t
test_slice_repeat()
{
    int_slice_repeat(s1, 88, 10);

    CC_ASSERT_INT_EQUAL(s1->len, 10);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, 0), 88);
    CC_ASSERT_INT_EQUAL(int_slice_get(s1, s1->len-1), 88);

    CC_SUCCESS;
}

cc_result_t
test_slice_grow()
{
    int_slice_grow(s1, 2);

    CC_ASSERT_UINT64_EQUAL(s1->cap, 5);

    CC_SUCCESS;
}

cc_result_t
test_slice_concat()
{
    int_slice_t *s2 = int_slice_new(3);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);

    CC_ASSERT_UINT64_EQUAL(s1->len, 3);
    CC_ASSERT_UINT64_EQUAL(s1->cap, 3);

    int_slice_append(s2, 9);
    int_slice_append(s2, 9);
    int_slice_append(s2, 9);

    CC_ASSERT_UINT64_EQUAL(s2->len, 3);
    CC_ASSERT_UINT64_EQUAL(s2->cap, 3);

    uint64_t s1_len = int_slice_concat(s1, s2);

    CC_ASSERT_UINT64_EQUAL(s1_len, 6); // purposefully wrong... correct (6)

    int_slice_free(s2);

    CC_SUCCESS;
}

int
main(int argc, char **argv)
{
    CC_INIT;

    CC_RUN(test_slice_append);
    CC_RUN(test_slice_get);
    CC_RUN(test_slice_contains);
    CC_RUN(test_slice_delete);
    CC_RUN(test_slice_replace_by_index);
    CC_RUN(test_slice_count);
    CC_RUN(test_slice_sort);
    CC_RUN(test_slice_replace_by_val);
    CC_RUN(test_slice_repeat);
    CC_RUN(test_slice_grow);
    CC_RUN(test_slice_concat);

    CC_COMPLETE;
}
