#include <windows.h>
#include "hotkey.h"

extern void handleHotkey(int id);

void setupHotkey() {
    RegisterHotKey(NULL, 1, MOD_ALT, 'R');
    RegisterHotKey(NULL, 2, MOD_ALT, 'Q');

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