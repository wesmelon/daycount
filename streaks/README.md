# Daycount streak tracker

A dependency-free, mobile-first app alongside the original Go application. No build, account, or database required. The legacy application is unchanged.

## Run

From the repository root: `python3 -m http.server 8080`, then open `http://localhost:8080/streaks/`. On an iPhone on the same network, use your computer's LAN IP instead of localhost. For home-screen installation and offline support, host the `streaks` directory on HTTPS (for example GitHub Pages), then use Safari → Share → Add to Home Screen. No hosting configuration is changed automatically by this feature.

If GitHub Pages serves the repository root, the app will be at `https://wesmelon.github.io/daycount/streaks/` after the branch is merged and Pages is enabled. This is an expected URL, not a claim that deployment is active.

## Features and data

- Custom categories, daily check-in/undo, current/longest streaks and total days.
- Monthly calendar with editable days and notes; historical inclusive date ranges.
- Badges at 10, 30, 50, 100 and 365 consecutive days, including backfilled records.
- Local browser storage, JSON export/validated restore, offline app shell.
- Backup/import transfers records between devices. There is **no server sync** or integration with legacy PostgreSQL accounts. Browser data can be cleared/evicted, so export regularly. Safari and a home-screen app may have separate storage: export before switching and import if necessary.
- Dates follow the device's local calendar; one entry per category/date. Yesterday's streak stays active until today is missed. Travel follows the device date. Future dates are rejected. Backfill overwrites notes on overlapping dates; the form explains this before saving.
- Badge awards derive from the longest saved streak. Editing/deleting dates can remove a badge. Badges are visual rewards, not financial bonuses.

## Test

`node --test streaks/model.test.mjs`

The code uses native ES modules and modern Safari APIs. Serve over HTTP/HTTPS; opening index.html directly as a file is unsupported. When changing offline assets, bump the cache version in `sw.js`. The legacy Go app is separate; its routes do not automatically expose this standalone app.
