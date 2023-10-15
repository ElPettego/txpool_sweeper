#include <thread>
#include <iostream>
#include <csignal>
#include <cstdlib>


int main(int argc, char *argv[]){    
    std::system("make --no-print-directory crun &");

    signal(SIGINT, [](int) {
        std::system("pkill -e -f ./src/scrape_eigen.py");
        std::system("pkill -e -f ./build/txpl_swpr");
        std::exit(0);
    });

    // Keep the main thread alive
    while (true) {
        int result = std::system("make --no-print-directory scrape");

        if (result) {
            std::system("pkill -e -f ./src/scrape_eigen.py");
            std::system("pkill -e -f ./build/txpl_swpr");
            std::exit(1);
        }

        // std::this_thread::sleep_for(std::chrono::seconds(70));
    }

    return 0;
}