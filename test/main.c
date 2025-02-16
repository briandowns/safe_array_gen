#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "int_slice.h"

typedef struct {
    char name[16];
    unsigned short age; 
    int_slice_t *grades;
} person_t;

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
main(int argc, char **argv)
{
    printf("Running tests...\n");
    person_t *p = (person_t*)malloc(sizeof(person_t));
    strcpy(p->name, "Brian");
    p->age = 44;
    p->grades = int_slice_new(4);
    assert(p->grades->cap == 4);
    assert(p->grades->len == 0);

    int_slice_append(p->grades, 9);
    int_slice_append(p->grades, 100);
    int_slice_append(p->grades, 7);
    int_slice_append(p->grades, 88);
    int_slice_append(p->grades, 70);
    int_slice_append(p->grades, 41);
    int_slice_append(p->grades, 15);
    int_slice_append(p->grades, 1);
    assert(p->grades->len == 8);
    assert(p->grades->cap == 8);
    int_slice_foreach(p->grades, print_item, NULL);

    assert(int_slice_contains(p->grades, 15));
    assert(!int_slice_contains(p->grades, 2));

    int_slice_delete(p->grades, 3);
    assert(p->grades->len == 7);

    int_slice_foreach(p->grades, print_item, NULL);

    assert(int_slice_get(p->grades, 1) == 100);

    int_slice_replace_by_idx(p->grades, 5, 42);
    assert(int_slice_get(p->grades, 5) == 42);
    int_slice_foreach(p->grades, print_item, NULL);

    int_slice_sort(p->grades);
    int_slice_foreach(p->grades, print_item, NULL);

    p->grades->compare = compare_func;
    int_slice_append(p->grades, 100);
    int_slice_append(p->grades, 100);
    int_slice_append(p->grades, 100);
    int_slice_append(p->grades, 100);
    int_slice_replace_by_val(p->grades, 100, 700, 3);
    int_slice_foreach(p->grades, print_item, NULL);

    int_slice_repeat(p->grades, 88, 10);
    assert(p->grades->len == 21);
    int_slice_foreach(p->grades, print_item, NULL);

    size_t count = int_slice_count(p->grades, 88);

    assert(count == 10);
    int_slice_grow(p->grades, 2);
    assert(p->grades->cap == 34);


    int_slice_t *s1 = int_slice_new(3);
    int_slice_t *s2 = int_slice_new(3);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    int_slice_append(s1, 7);
    printf("%lu\n", s1->len);
    assert(s1->len == 3);
    assert(s1->cap == 3);
    int_slice_append(s2, 9);
    int_slice_append(s2, 9);
    int_slice_append(s2, 9);
    assert(s2->len == 3);
    assert(s2->cap == 3);

    size_t s1_len = int_slice_concat(s1, s2);
    printf("The new array: %lu\n", s1_len);
    int_slice_foreach(s1, print_item, NULL);

    int_slice_free(p->grades);
    int_slice_free(s1);
    int_slice_free(s2);
    free(p);

    printf("Tests ran successfully!...\n");

    return 0;
}
