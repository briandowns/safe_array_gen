#include <assert.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>

#include "uint8_t_slice.h"

typedef struct {
    char name[16];
    unsigned short age; 
    uint8_slice_t *grades;
} person_t;

void
print_item(uint8_t item, void *user_data)
{
    printf("%d\n", item);
}

int
main(int argc, char **argv)
{
    printf("Running tests...\n");
    person_t *p = (person_t*)malloc(sizeof(person_t));
    strcpy(p->name, "Brian");
    p->age = 44;
    p->grades = uint8_slice_new(4);
    assert(p->grades->cap == 4);
    assert(p->grades->len == 0);

    uint8_slice_append(p->grades, 9);
    uint8_slice_append(p->grades, 100);
    uint8_slice_append(p->grades, 7);
    uint8_slice_append(p->grades, 88);
    uint8_slice_append(p->grades, 70);
    uint8_slice_append(p->grades, 41);
    uint8_slice_append(p->grades, 15);
    uint8_slice_append(p->grades, 1);
    assert(p->grades->len == 8);
    assert(p->grades->cap == 8);

    assert(uint8_slice_contains(p->grades, 15));
    assert(!uint8_slice_contains(p->grades, 2));

    uint8_slice_delete(p->grades, 3);
    assert(p->grades->len == 7);

    assert(uint8_slice_get(p->grades, 1) == 100);

    uint8_slice_replace(p->grades, 5, 42);
    assert(uint8_slice_get(p->grades, 5) == 42);

    uint8_slice_foreach(p->grades, print_item, NULL);

    uint8_slice_free(p->grades);
    free(p);

    printf("Tests ran successfully!...\n");

    return 0;
}
