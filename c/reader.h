#ifndef _MAL_READER_H_  // NOLINT
#define _MAL_READER_H_

typedef struct reader {
    char *line;
} Reader;

Reader *Reader_new();
char *Reader_peek(Reader *r);
char *Reader_next(Reader *r);
void Reader_free(Reader *r);

#endif // _MAL_READER_H_
