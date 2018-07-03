#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "types.h"

void MalString_free(MalString *s);
void MalList_free(MalList *l);
void MalMap_free(MalMap *l);

MalValue *MalValue_new(MalType type)
{
    MalValue *v = (MalValue*)malloc(sizeof(MalValue));
    if (v == NULL) {
        perror("OOM");
        exit(1);
    }
    v->type = type;
    return v;
}

void MalValue_free(MalValue *v)
{
    switch (v->type) {
        case MAL_STRING:
        case MAL_SYMBOL:
            MalString_free(v->u.str);
            v->u.str = NULL;
            break;
        case MAL_LIST:
        case MAL_VECTOR:
            MalList_free(v->u.list);
            v->u.list = NULL;
            break;
        case MAL_MAP:
            MalMap_free(v->u.map);
            v->u.map = NULL;
            break;
        case MAL_FUNC:
            break;
        default: break;
    }
    free(v);
}

MalValue *MalValue_new_nil()
{
    MalValue *v = MalValue_new(MAL_NIL);
    return v;
}

MalValue *MalValue_new_int(int n)
{
    MalValue *v = MalValue_new(MAL_INTEGER);
    v->u.i = n;
    return v;
}

MalValue *MalValue_new_float(double n)
{
    MalValue *v = MalValue_new(MAL_FLOAT);
    v->u.f = n;
    return v;
}

MalValue *MalValue_new_string(char *s, size_t len)
{
    MalValue *v = MalValue_new(MAL_STRING);
    v->u.str = (MalString*)malloc(sizeof(MalString));
    v->u.str->s = (char*)malloc(sizeof(char)*len);
    memcpy(v->u.str->s, s, len);
    v->u.str->len = len;
    return v;
}

MalValue *MalValue_new_symbol(char *s, size_t len)
{
    MalValue *v = MalValue_new_string(s, len);
    v->type = MAL_SYMBOL;
    return v;
}

MalValue *MalValue_new_list()
{
    MalValue *v = MalValue_new(MAL_LIST);
    v->u.list = NULL;  // TODO(damnever)
    return v;
}

MalValue *MalValue_new_vector()
{
    MalValue *v = MalValue_new_list();
    v->type = MAL_VECTOR;
    return v;
}

MalValue *MalValue_new_map()
{
    MalValue *v = MalValue_new(MAL_MAP);
    v->u.map = NULL;  // TODO(damnever);
    return v;
}
