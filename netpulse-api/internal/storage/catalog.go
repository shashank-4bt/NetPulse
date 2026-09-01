package storage

import "github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"

func ServiceCatalog() []contract.Service {
	return []contract.Service{
		{Slug: "google", Name: "Google", Category: "Search & identity", Summary: "Search, accounts, and related Google properties.", Layers: []string{"DNS", "Routing", "CDN", "TLS", "HTTP", "Service"}},
		{Slug: "youtube", Name: "YouTube", Category: "Media", Summary: "Video delivery and playback endpoints.", Layers: []string{"DNS", "CDN", "TLS", "HTTP", "Service"}},
		{Slug: "cloudflare", Name: "Cloudflare", Category: "Infrastructure", Summary: "DNS, CDN, and edge connectivity used by many sites.", Layers: []string{"DNS", "Connectivity", "Routing", "CDN", "TLS"}},
		{Slug: "github", Name: "GitHub", Category: "Developer", Summary: "Git hosting, APIs, and web application.", Layers: []string{"DNS", "TLS", "HTTP", "Service"}},
		{Slug: "microsoft-365", Name: "Microsoft 365", Category: "Productivity", Summary: "Identity, mail, and collaboration endpoints.", Layers: []string{"DNS", "ISP", "TLS", "HTTP", "Service"}},
		{Slug: "slack", Name: "Slack", Category: "Collaboration", Summary: "Messaging clients and real-time service endpoints.", Layers: []string{"DNS", "CDN", "TLS", "HTTP", "Service"}},
		{Slug: "aws", Name: "AWS", Category: "Cloud", Summary: "Regional control-plane and commonly used public endpoints.", Layers: []string{"DNS", "Routing", "TLS", "HTTP", "Service"}},
		{Slug: "zoom", Name: "Zoom", Category: "Meetings", Summary: "Meeting join paths and media edge connectivity.", Layers: []string{"Connectivity", "ISP", "Routing", "TLS", "Service"}},
		{Slug: "instagram", Name: "Instagram", Category: "Social", Summary: "App and web media endpoints.", Layers: []string{"DNS", "CDN", "TLS", "HTTP", "Service"}},
	}
}
