#include <stdbool.h>
#include <stdint.h>

void petInit(void);
int petScreen(int i, double *x, double *y, double *w, double *h, double *scale);
bool petCreate(double x, double y, double w, double h, double scale, bool topmost,
               const uint8_t *pix, int frames, int pw, int ph);
void petShowFrame(int i);
void petRun(double interval);
