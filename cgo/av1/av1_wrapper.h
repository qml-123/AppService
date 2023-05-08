// av1_wrapper.h

#ifndef AV1_WRAPPER_H
#define AV1_WRAPPER_H

#include <stdint.h>
#include <stddef.h>

int av1_encode(const uint8_t *input, size_t input_size, uint8_t **output, size_t *output_size);
int av1_decode(const uint8_t *input, size_t input_size, uint8_t **output, size_t *output_size);

#endif // AV1_WRAPPER_H