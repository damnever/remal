#ifndef _MAL_TYPES_H_ // NOLINT
#define _MAL_TYPES_H_

#include <stddef.h>

#define DEFAULT_CONTAINER_CAP 1

typedef enum mal_type {
    MAL_NIL,
    MAL_TRUE,
    MAL_FALSE,
    MAL_INTEGER,
    MAL_FLOAT,
    MAL_STRING,
    MAL_SYMBOL,
    MAL_LIST,
    MAL_VECTOR,
    MAL_MAP,
    MAL_FUNC,
} MalType;

typedef struct mal_value MalValue;

typedef struct mal_string {
    char *s;  // null-terminated
    size_t len;
} MalString;

typedef struct mal_list_elem {
    MalValue *val;
    struct mal_list_elem *next;
} MalListElem;

typedef struct mal_list {
    size_t size;
    MalListElem *elems;
} MalList;

typedef struct mal_map_elem {
    MalString *key;
    MalValue *val;
    struct mal_map_elem *next;
} MalMapElem;

typedef struct mal_map {
    size_t size;
    size_t cap;
    MalMapElem **elems;
} MalMap;

typedef struct mal_func {
} MalFunc;

typedef struct mal_value {
    MalType type;
    union {
        int i;  // integer
        double f;  // float
        MalString *str;  // string or symbol
        MalList *list; //
        MalMap *map;
        MalFunc *func;
    } u;
} MalValue;

MalValue *MalValue_new_nil();
MalValue *MalValue_new_int(int n);
MalValue *MalValue_new_float(double n);
MalValue *MalValue_new_string(char *s, size_t len);
MalValue *MalValue_new_symbol(char *s, size_t len);
MalValue *MalValue_new_list();
MalValue *MalValue_new_vector();
MalValue *MalValue_new_map();
MalValue *MalValue_new(MalType type);
void MalValue_free(MalValue *v);

#endif
