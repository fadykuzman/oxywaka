# OxyWaka Licensing Notice

This repository (oxywaka/wakapi-core) contains the **open-source core** of OxyWaka.

---

## 📜 Two-Tier License Structure

### 1. OxyWaka Core (THIS REPOSITORY) - MIT License

**What's included (Free & Open Source):**
- ✅ WakaTime-compatible heartbeat tracking
- ✅ Time summaries and statistics
- ✅ User authentication and management
- ✅ Team leaderboards
- ✅ Basic dashboards
- ✅ REST API
- ✅ Database schema and migrations

**License:** MIT (see [LICENSE](./LICENSE))

**Based on:** [Wakapi](https://github.com/muety/wakapi) by Ferdinand Mütsch

You are free to:
- ✅ Use commercially
- ✅ Modify
- ✅ Distribute
- ✅ Sublicense
- ✅ Use privately

**Requirements:**
- Include original MIT license
- Include copyright notice

---

### 2. OxyWaka Insights (Separate Repository) - Proprietary

**What's NOT included here (Commercial License Required):**
- ❌ AI productivity scoring algorithms
- ❌ Deep work detection engine
- ❌ Advanced pattern analysis
- ❌ Context switch quantification
- ❌ AI-powered recommendations
- ❌ Burnout risk detection
- ❌ Performance forecasting

**Repository:** [Private - Licensed customers only]

**License:** Proprietary commercial license

**Available via:**
- SaaS subscription: https://oxywaka.com ($8-15/user/month)
- Enterprise self-hosted: Contact enterprise@oxywaka.com ($15k+/year)
- Commercial reseller: Contact partnerships@oxywaka.com

---

## 🏗️ How They Work Together

```
┌─────────────────────────────────────────┐
│  OxyWaka Core (MIT - This Repo)         │
│  ┌──────────────────────────────────┐   │
│  │ • Heartbeat collection           │   │
│  │ • Time tracking                  │   │
│  │ • Basic summaries                │   │
│  │ • User management                │   │
│  │ • API endpoints                  │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
                  ↕ REST API
┌─────────────────────────────────────────┐
│  OxyWaka Insights (Proprietary)          │
│  ┌──────────────────────────────────┐   │
│  │ • AI productivity analysis       │   │
│  │ • Deep work detection            │   │
│  │ • Pattern recognition            │   │
│  │ • Smart recommendations          │   │
│  └──────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

The insights engine consumes data from the core via standard REST APIs.
No proprietary code is mixed with the MIT-licensed core.

---

## 🤔 Which Should I Use?

### Use OxyWaka Core (Free/MIT) if:
- ✅ You want basic time tracking
- ✅ You need WakaTime compatibility
- ✅ You can self-host
- ✅ You don't need AI insights
- ✅ You want to contribute to open source

### Use OxyWaka SaaS (Paid) if:
- ✅ You want AI-powered productivity insights
- ✅ You prefer managed hosting (no DevOps)
- ✅ You need team analytics
- ✅ You want priority support
- ✅ You value time over cost

### Use OxyWaka Enterprise (Paid) if:
- ✅ You need on-premise deployment
- ✅ You have compliance requirements
- ✅ You want AI insights on your infrastructure
- ✅ You need SLA and dedicated support
- ✅ You have 100+ developers

---

## 🔍 Comparison

| Feature | OxyWaka Core (Free) | OxyWaka SaaS | OxyWaka Enterprise |
|---------|---------------------|--------------|-------------------|
| **Time tracking** | ✅ | ✅ | ✅ |
| **Basic dashboards** | ✅ | ✅ | ✅ |
| **Team features** | ✅ | ✅ | ✅ |
| **Self-hosted** | ✅ | ❌ | ✅ |
| **AI insights** | ❌ | ✅ | ✅ |
| **Deep work detection** | ❌ | ✅ | ✅ |
| **Managed hosting** | ❌ | ✅ | Optional |
| **Priority support** | ❌ | ✅ | ✅ |
| **SLA** | ❌ | ⚠️ 99% | ✅ 99.9% |
| **Price** | Free | $8-15/user/mo | $15k+/year |

---

## 💡 For Developers

### Contributing to Core
We welcome contributions to the open-source core!
- Report bugs
- Submit PRs
- Suggest features
- Improve documentation

See [CONTRIBUTING.md](./CONTRIBUTING.md)

### Building Your Own Insights
Want to build proprietary features on top?
- ✅ Fork this repo (MIT license allows it)
- ✅ Build your own analytics
- ✅ Keep this license notice
- ✅ You don't have to open-source your additions

### Using the API
The core exposes REST APIs for:
- Heartbeat ingestion
- Summary retrieval
- User management
- Duration queries

Build your own tools on top! (MIT allows commercial use)

---

## 📞 Questions?

**Open Source (Core):**
- Issues: https://github.com/oxywaka/wakapi-core/issues
- Discussions: https://github.com/oxywaka/wakapi-core/discussions

**Commercial (Insights):**
- Sales: sales@oxywaka.com
- Support: support@oxywaka.com
- Partnerships: partnerships@oxywaka.com

**General:**
- Website: https://oxywaka.com
- Email: hello@oxywaka.com

---

## ⚖️ Legal

**OxyWaka Core:** Copyright © 2025 OxyWaka (MIT License)
**Original Wakapi:** Copyright © 2020 Ferdinand Mütsch (MIT License)
**OxyWaka Insights:** Copyright © 2025 OxyWaka (Proprietary License)

**Trademarks:** OxyWaka™ is a trademark of OxyWaka.

---

**Built with ❤️ for developers.**

*Open-source core. Proprietary AI. Your choice.*
