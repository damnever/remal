#include <stdio.h>
#include <stdlib.h>

#include <readline/history.h>
#include <readline/readline.h>

#define ERR_NULL


char *read()
{
    char *line = readline("> ");
    if (line != NULL) {
        add_history(line);
    }
    return line;
}

char *eval(char *line)
{
    return line;
}

void print(char *line)
{
    puts(line);
}

int rep()
{
    char *line = read();
    if (line == NULL) { return -1; }

    print(eval(line));
    free(line);
    return 0;
}

int main(void)
{
    int r = 0;
    puts("Have fun!");

    for (;;) {
        if ((r = rep()) != 0) {
            break;
        }
    }

    puts("\nBye.");
    return 0;
}
