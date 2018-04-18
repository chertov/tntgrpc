#pragma once
#include <string>

bool run_server(const std::string server_address, bool sync_mode);
void grpc_start(const std::string server_address, bool sync_mode);
