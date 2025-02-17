#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "int_slice.h"
#include "utest.h"

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

bool
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

    CC_ASSERT_EQUAL_SIZE_T(s1->len, (size_t)8);
    CC_ASSERT_EQUAL_SIZE_T(s1->cap, (size_t)8);

    int_slice_free(s1);

    return true;
}

bool
test_slice_get()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);

    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 0), 9);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 1), 100);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 2), 7);

    int_slice_free(s1);

    return true;
}

bool
test_slice_contains()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    
    CC_ASSERT(int_slice_contains(s1, 9));
    CC_ASSERT(int_slice_contains(s1, 100));
    CC_ASSERT(int_slice_contains(s1, 7));

    int_slice_free(s1);

    return true;
}

bool
test_slice_delete()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    
    int_slice_delete(s1, 2);

    CC_ASSERT_EQUAL_SIZE_T(s1->len, 2);

    int_slice_free(s1);

    return true;
}

bool
test_slice_replace_by_index()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_replace_by_idx(s1, 0, 42);

    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 0), 42);

    int_slice_free(s1);

    return true;
}

bool
test_slice_count()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 9);
    int_slice_append(s1, 100);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);

    size_t count = int_slice_count(s1, 7);
    CC_ASSERT_EQUAL_SIZE_T(count, 3);

    int_slice_free(s1);

    return true;
}

bool
test_slice_sort()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_append(s1, 100);
    int_slice_append(s1, 1000);
    int_slice_append(s1, 9);

    int_slice_sort(s1);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 0), 9);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 1), 100);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 2), 1000);

    int_slice_free(s1);

    return true;
}

bool
test_slice_replace_by_val()
{
    int_slice_t *s1 = int_slice_new(3);
    s1->compare = compare_func;

    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);
    int_slice_append(s1, 100);

    int_slice_replace_by_val(s1, 100, 700, 3);

    CC_ASSERT_EQUAL_SIZE_T(int_slice_count(s1, 700), 3);
    CC_ASSERT_EQUAL_SIZE_T(int_slice_count(s1, 100), 1);

    int_slice_free(s1);

    return true;
}

bool
test_slice_repeat()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_repeat(s1, 88, 10);

    CC_ASSERT_EQUAL_INT(s1->len, 10);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, 0), 88);
    CC_ASSERT_EQUAL_INT(int_slice_get(s1, s1->len-1), 88);

    int_slice_free(s1);

    return true;
}

bool
test_slice_grow()
{
    int_slice_t *s1 = int_slice_new(3);

    int_slice_grow(s1, 2);
    CC_ASSERT_EQUAL_SIZE_T(s1->cap, 5);

    int_slice_free(s1);

    return true;
}

bool
test_slice_concat()
{
    int_slice_t *s1 = int_slice_new(3);
    int_slice_t *s2 = int_slice_new(3);

    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);

    CC_ASSERT_EQUAL_SIZE_T(s1->len, 3);
    CC_ASSERT_EQUAL_SIZE_T(s1->cap, 3);

    int_slice_append(s2, 9);
    int_slice_append(s2, 9);
    int_slice_append(s2, 9);

    CC_ASSERT_EQUAL_SIZE_T(s2->len, 3);
    CC_ASSERT_EQUAL_SIZE_T(s2->cap, 3);

    size_t s1_len = int_slice_concat(s1, s2);
    CC_ASSERT_EQUAL_SIZE_T(s1_len, 5); // purposefully wrong... correct (6)

    int_slice_free(s1);
    int_slice_free(s2);

    return true;
}

int
main(int argc, char **argv)
{
    utest_init();

    CC_RUN("append", test_slice_append);
    CC_RUN("get", test_slice_get);
    CC_RUN("contains", test_slice_contains);
    CC_RUN("delete", test_slice_delete);
    CC_RUN("replace by index", test_slice_replace_by_index);
    CC_RUN("count", test_slice_count);
    CC_RUN("sort", test_slice_sort);
    CC_RUN("replace by value", test_slice_replace_by_val);
    CC_RUN("repeat", test_slice_repeat);
    CC_RUN("grow", test_slice_grow);
    CC_RUN("concat", test_slice_concat);

    utest_complete();

    return 0;
}
