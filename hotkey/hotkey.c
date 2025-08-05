// hotkey/hotkey.c
#include <windows.h>
#include <stdio.h>

static int running = 1;

// Esta función bloqueante registra y espera el hotkey Ctrl+Alt+B (0x42)
void ListenHotkey() {
    if (!RegisterHotKey(NULL, 1, MOD_CONTROL | MOD_ALT, 0x42)) {
        printf("Failed to register hotkey\n");
        return;
    }
    printf("Hotkey registered: Ctrl+Alt+B\n");

    MSG msg;
    while (running && GetMessage(&msg, NULL, 0, 0)) {
        if (msg.message == WM_HOTKEY) {
            printf("Hotkey pressed!\n");
            // Aquí podrías enviar señal a Go o hacer algo
        }
    }

    UnregisterHotKey(NULL, 1);
}

void StopHotkey() {
    running = 0;
    // Se podría enviar un mensaje para salir del GetMessage, 
    // pero para simplificar no lo hacemos ahora
}
