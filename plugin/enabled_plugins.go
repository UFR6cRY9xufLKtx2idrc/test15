/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package plugin

// data providers
import (
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/data_provider/domain_set"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/data_provider/ip_set"
)

// matches
import (
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/client_ip"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/cname"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/env"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/has_resp"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/has_wanted_ans"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/ptr_ip"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/qclass"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/qname"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/qtype"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/random"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/rcode"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/matcher/resp_ip"
)

// executables
import (
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/arbitrary"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/black_hole"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/cache"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/debug_print"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/drop_resp"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/dual_selector"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/ecs"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/forward"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/hosts"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/ipset"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/metrics_collector"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/nftset"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/query_summary"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/redirect"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/reverse_lookup"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/sequence"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/sequence/fallback"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/sleep"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/executable/ttl"
)

// other
import (
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/mark" // executable and matcher
)

// servers
import (
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/server/http_server"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/server/tcp_server"
	_ "github.com/UFR6cRY9xufLKtx2idrc/mosdns/main/plugin/server/udp_server"
)
