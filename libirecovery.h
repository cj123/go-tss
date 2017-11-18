/*
 * libirecovery.h stub, used to allow inclusion of idevicerestore/common.h without modifications.
 *
 * full file can be found at: https://github.com/libimobiledevice/libirecovery/blob/master/include/libirecovery.h
 */

struct irecv_device {
    const char* product_type;
    const char* hardware_model;
    unsigned int board_id;
    unsigned int chip_id;
};
typedef struct irecv_device* irecv_device_t;
