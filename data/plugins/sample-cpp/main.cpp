#include <iostream>
#include <fstream>
#include <cstdlib>
#include <thread>
#include <chrono>
#include <sstream>

int main(int argc, char* argv[]) {
    // These values are created by Achaean
    const char* pluginID = std::getenv("PLUGIN_ID");
    const char* progressPipeFile = std::getenv("PROGRESS_PIPE");

    std::cout << "Hello from plugin " << pluginID << std::endl;

    // Convert args in a string
    std::ostringstream paramsStream;
    for (int i = 1; i < argc; ++i) {
        if (i > 1) paramsStream << " ";
        paramsStream << argv[i];
    }
    std::cout << "Parameters: " << paramsStream.str() << std::endl;

    // Open Named Pipe in RO mode
    std::ofstream progressPipe(progressPipeFile);
    if (!progressPipe) {
        std::cerr << "Error opening named pipe: " << progressPipeFile << std::endl;
        return 1;
    }

    // Dummy progress. 0% to 100% (5% increment).
    for (int p = 0; p <= 100; p += 5) {
        if (p == 50) {
            std::cout << "Half of the progress..." << std::endl;  // Stdout
            std::cerr << "This is a fake error." << std::endl;    // Stderr
        }
        if (p == 100) {
            std::cout << "Hasta la vista baby." << std::endl;
        }

        progressPipe << p << std::endl;
        progressPipe.flush(); // Force flush to empty buffer.
        std::this_thread::sleep_for(std::chrono::milliseconds(250));
    }

    // IMPORTANT: Plugins must end writing "DONE" to the progressPipe before exit.
    progressPipe << "DONE" << std::endl;
    progressPipe.flush();

    return 0;
}
