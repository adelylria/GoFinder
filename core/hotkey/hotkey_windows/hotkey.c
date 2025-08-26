#ifdef _WIN32
#include <windows.h>
#include "hotkey.h"

// Declaración de la función Go (se define en Go via //export)
extern void handleHotkey(int id);

void setupHotkey() {
    RegisterHotKey(NULL, 1, MOD_ALT, 0x52); // 0x52 is 'R'
    RegisterHotKey(NULL, 2, MOD_ALT, 0x51); // 0x51 is 'Q'

    MSG msg = {0};
    while (GetMessage(&msg, NULL, 0, 0) != 0) {
        if (msg.message == WM_HOTKEY) {
            handleHotkey((int)msg.wParam);
            if (msg.wParam == 2) {
                break;
            }
        }
    }

    UnregisterHotKey(NULL, 1);
    UnregisterHotKey(NULL, 2);
}

#endif