import { useEffect, useState, type ReactNode } from "react";
import {
  benchmarkMetrics,
  downloadUrl,
  faqs,
  githubUrl,
  loadLatestRelease,
  platformDownloads,
  problemItems,
  protocolFields,
  providers,
  releaseLoadingVersion,
  workflowSteps,
  type LatestRelease,
} from "../lib/content/site";
import { ReleaseVersion } from "./shared";

export function Problem() {
  return (
    <section className="border-b border-[var(--hawp-line)] py-24 sm:py-28">
      <div className="mx-auto grid max-w-[1280px] gap-14 px-5 sm:px-8 lg:grid-cols-[.92fr_1.08fr] lg:px-10">
        <div>
          <Eyebrow>The failure mode</Eyebrow>
          <h2 className="mt-5 max-w-xl text-4xl font-bold leading-[1.08] tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
            AI work often loses shape between the idea and execution.
          </h2>
          <p className="mt-6 max-w-xl text-lg leading-8 text-[var(--hawp-muted)]">
            The expensive part is not always generation. It is re-explaining
            intent, rediscovering constraints, recovering decisions, and
            correcting drift across sessions and tools.
          </p>
        </div>
        <div className="divide-y divide-[var(--hawp-line)] overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)]">
          {problemItems.map((item) => (
            <article
              key={item[0]}
              className="grid grid-cols-[46px_1fr] gap-4 p-6 sm:p-7"
            >
              <span className="font-mono text-sm text-[var(--hawp-clay)]">
                {item[0]}
              </span>
              <div>
                <h3 className="text-xl font-semibold text-[var(--hawp-ink)]">
                  {item[1]}
                </h3>
                <p className="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">
                  {item[2]}
                </p>
              </div>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
export function Eyebrow({ children }: { children: ReactNode }) {
  return (
    <p className="text-xs font-bold tracking-[0.08em] text-[var(--hawp-clay)] uppercase">
      {children}
    </p>
  );
}
export function Protocol() {
  return (
    <section
      id="protocol"
      className="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
        <Eyebrow>The protocol</Eyebrow>
        <h2 className="mt-5 max-w-3xl text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
          Five fields. One portable shape for the work.
        </h2>
        <p className="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">
          HAWP keeps the protocol intentionally small. The artifact remains
          readable by humans, useful to agents, and portable across tools.
        </p>
        <div className="mt-12 grid gap-3 md:grid-cols-2 lg:grid-cols-5">
          {protocolFields.map((f) => (
            <article
              key={f[0]}
              className="min-h-56 rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-6"
            >
              <span className="font-mono text-xs text-[var(--hawp-clay)]">
                {f[0]}
              </span>
              <h3 className="mt-8 text-xl font-semibold text-[var(--hawp-ink)]">
                {f[1]}
              </h3>
              <p className="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">
                {f[2]}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
export function Agents({
  version,
  loading,
}: {
  version: string;
  loading: boolean;
}) {
  return (
    <section
      id="agents"
      className="border-b border-[var(--hawp-line)] py-24 sm:py-28"
    >
      <div className="mx-auto grid max-w-[1280px] items-center gap-14 px-5 sm:px-8 lg:grid-cols-2 lg:px-10">
        <div className="overflow-hidden rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-alt)] shadow-2xl shadow-black/20">
          <div className="flex gap-2 border-b border-[var(--hawp-line)] px-5 py-4">
            {[0, 1, 2].map((n) => (
              <i
                key={n}
                className="size-2.5 rounded-full bg-[var(--hawp-red)]"
              />
            ))}
          </div>
          <div className="p-6 font-mono text-sm leading-7 text-[var(--hawp-muted)] sm:p-8">
            <div>
              <span className="text-[var(--hawp-clay)]">$</span> hawp work new
            </div>
            <div className="text-[var(--hawp-clay)]">
              ✓ work item scaffolded
            </div>
            {[
              "input       → captured",
              "context     → bounded",
              "mission     → explicit",
              "constraints → preserved",
              "output      → verifiable",
            ].map((line) => (
              <div key={line} className="text-[var(--hawp-muted-2)]">
                &nbsp;&nbsp;{line}
              </div>
            ))}
            <div className="mt-5">
              <span className="text-[var(--hawp-clay)]">$</span> hawp mcp
            </div>
            <div className="text-[var(--hawp-clay)]">
              ✓ hawp_work_intake · compound shaping
            </div>
            <div className="text-[var(--hawp-clay)]">
              ✓ hawp_work_doc · canonical work paths
            </div>
          </div>
        </div>
        <div>
          <Eyebrow>Tool-independent by design</Eyebrow>
          <h2 className="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
            The work survives the agent.
          </h2>
          <p className="mt-5 text-lg leading-8 text-[var(--hawp-muted)]">
            HAWP uses Markdown artifacts and exposes search and workflows
            through MCP, so the same project shape can move across coding agents
            without getting trapped in one vendor&apos;s memory or chat history.
          </p>
          <p className="mt-4 text-sm leading-6 text-[var(--hawp-muted-2)]">
            In <ReleaseVersion version={version} loading={loading} />, compound
            intake can retrieve relevant context, reshape the request through a
            local LLM, and return structured mission, constraints, output
            specification, and done signal states without creating a work item
            automatically.
          </p>
          <div className="mt-7 flex flex-wrap gap-2">
            {providers.map((provider) => (
              <span
                key={provider}
                className="rounded-full border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-3 py-1.5 text-xs text-[var(--hawp-muted)]"
              >
                {provider}
              </span>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
export function Benchmarks({
  version,
  loading,
}: {
  version: string;
  loading: boolean;
}) {
  return (
    <section
      id="benchmarks"
      className="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
        <Eyebrow>
          <ReleaseVersion version={version} loading={loading} /> repository
          evidence
        </Eyebrow>
        <h2 className="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
          Less context. Faster retrieval. Explicit uncertainty.
        </h2>
        <p className="mt-5 max-w-3xl text-lg leading-8 text-[var(--hawp-muted)]">
          Current repository evidence covers search retrieval, context shaping,
          and structured intake shaping. These are HAWP-local measurements
          intended to make workflow claims inspectable.
        </p>
        <div className="mt-12 grid gap-3 md:grid-cols-3">
          {benchmarkMetrics.map((metric) => (
            <article
              key={metric[0]}
              className="rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-7"
            >
              <strong className="block text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-peach)]">
                {metric[0]}
              </strong>
              <span className="mt-4 block text-sm leading-6 text-[var(--hawp-muted-2)]">
                {metric[1]}
              </span>
            </article>
          ))}
        </div>
        <p className="mt-5 text-xs text-[var(--hawp-muted-2)]">
          Benchmarks are project-local measurements, not universal model
          performance claims.
        </p>
      </div>
    </section>
  );
}
export function Workflow() {
  return (
    <section className="border-b border-[var(--hawp-line)] py-24 sm:py-28">
      <div className="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
        <Eyebrow>How it fits</Eyebrow>
        <h2 className="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
          Shape, execute, and hand off.
        </h2>
        <div className="mt-12 grid gap-8 md:grid-cols-3">
          {workflowSteps.map((step) => (
            <article
              key={step[0]}
              className="border-l border-[var(--hawp-line-strong)] pl-5"
            >
              <span className="font-mono text-xs text-[var(--hawp-clay)]">
                {step[0]}
              </span>
              <h3 className="mt-5 text-xl font-semibold text-[var(--hawp-ink)]">
                {step[1]}
              </h3>
              <p className="mt-3 text-sm leading-6 text-[var(--hawp-muted-2)]">
                {step[2]}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
export function Downloads() {
  const [release, setRelease] = useState<LatestRelease | null>(null);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    loadLatestRelease()
      .then(setRelease)
      .finally(() => setLoading(false));
  }, []);
  const version = release?.tag_name.replace(/^v/, "") ?? releaseLoadingVersion;
  const asset = (name: string) =>
    release?.assets.find((candidate) => candidate.name === name)
      ?.browser_download_url ?? downloadUrl(name);
  return (
    <section
      id="downloads"
      className="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-[1280px] px-5 sm:px-8 lg:px-10">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <Eyebrow>Latest release</Eyebrow>
            <h2 className="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
              Get HAWP for your platform.
            </h2>
            <p className="mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">
              Download the HAWP CLI from the latest GitHub release. The version
              and asset links update from the release feed when a new build is
              published.
            </p>
          </div>
          <a
            className="shrink-0 text-sm font-semibold text-[var(--hawp-clay)] underline decoration-[var(--hawp-line-strong)] underline-offset-4 hover:text-[var(--hawp-ink)]"
            href={release?.html_url ?? `${githubUrl}/releases`}
          >
            View release notes ↗
          </a>
        </div>
        <p className="sr-only" aria-live="polite">
          {loading
            ? "Checking GitHub for the latest HAWP release."
            : `Showing HAWP version ${version}.`}
        </p>
        <div className="mt-10 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {platformDownloads.map((platform) => (
            <a
              key={platform.asset}
              className="flex items-center justify-between rounded-2xl border border-[var(--hawp-line)] bg-[var(--hawp-raised)] p-5 transition hover:border-[var(--hawp-red)] hover:bg-[var(--hawp-panel)]"
              href={asset(platform.asset)}
            >
              <span className="flex items-center gap-3">
                <span
                  className="grid size-9 place-items-center rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] text-sm text-[var(--hawp-clay)]"
                  aria-hidden="true"
                >
                  {platform.icon}
                </span>
                <span>
                  <strong className="block text-sm text-[var(--hawp-ink)]">
                    {platform.label}
                  </strong>
                  <span className="mt-1 block text-xs text-[var(--hawp-muted-2)]">
                    <ReleaseVersion version={version} loading={loading} />
                  </span>
                </span>
              </span>
              <span
                className="text-lg text-[var(--hawp-clay)]"
                aria-hidden="true"
              >
                ↓
              </span>
            </a>
          ))}
        </div>
      </div>
    </section>
  );
}
export function Faq() {
  return (
    <section
      id="faq"
      className="border-b border-[var(--hawp-line)] bg-[var(--hawp-alt)] py-24 sm:py-28"
    >
      <div className="mx-auto max-w-4xl px-5 sm:px-8">
        <Eyebrow>FAQ</Eyebrow>
        <h2 className="mt-5 text-4xl font-bold tracking-[-0.035em] text-[var(--hawp-ink)] sm:text-5xl">
          Small protocol, deliberate boundaries.
        </h2>
        <div className="mt-10 divide-y divide-[var(--hawp-line)] border-y border-[var(--hawp-line)]">
          {faqs.map((faq) => (
            <article key={faq[0]} className="py-6">
              <h3 className="text-lg font-semibold text-[var(--hawp-ink)]">
                {faq[0]}
              </h3>
              <p className="mt-2 text-sm leading-6 text-[var(--hawp-muted-2)]">
                {faq[1]}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
export function FinalCta() {
  return (
    <section className="py-24 sm:py-28">
      <div className="mx-auto max-w-[1180px] px-5 sm:px-8">
        <div className="site-grid overflow-hidden rounded-[2rem] border border-[var(--hawp-line-strong)] bg-[var(--hawp-raised)] px-6 py-12 text-center shadow-2xl shadow-black/20 sm:px-12 sm:py-16">
          <h2 className="mx-auto max-w-3xl text-4xl font-bold tracking-[-0.04em] text-[var(--hawp-ink)] sm:text-5xl">
            Give AI agent work a shape that survives the session.
          </h2>
          <p className="mx-auto mt-5 max-w-2xl text-lg leading-8 text-[var(--hawp-muted)]">
            HAWP is open source, Apache-2.0 licensed, and designed to stay
            readable even when the agent changes.
          </p>
          <div className="mt-8 flex flex-wrap justify-center gap-3">
            <a
              href={githubUrl}
              className="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]"
            >
              Read the repository ↗
            </a>
            <a
              href={`${githubUrl}#install--pick-your-agent`}
              className="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)]"
            >
              Get started
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
export function Footer() {
  return (
    <footer className="border-t border-[var(--hawp-line)] py-10">
      <div className="mx-auto flex max-w-[1280px] flex-col gap-5 px-5 text-sm text-[var(--hawp-muted-2)] sm:px-8 md:flex-row md:items-center md:justify-between lg:px-10">
        <span>HAWP · Human-AI Workflow Protocol · Sentzunhat</span>
        <div className="flex gap-5">
          <a className="hover:text-[var(--hawp-ink)]" href={githubUrl}>
            GitHub
          </a>
          <a
            className="hover:text-[var(--hawp-ink)]"
            href={`${githubUrl}/blob/main/LICENSE`}
          >
            Apache-2.0
          </a>
          <a
            className="hover:text-[var(--hawp-ink)]"
            href={`${githubUrl}/discussions`}
          >
            Discussions
          </a>
        </div>
      </div>
    </footer>
  );
}
export function AnalyticsConsent() {
  const [choice, setChoice] = useState<"unknown" | "accepted" | "declined">(
    "unknown",
  );
  const [visible, setVisible] = useState(false);
  const measurementId = "G-1FHPSYNPBS";
  const load = () => {
    if (
      document.querySelector(`script[data-hawp-analytics="${measurementId}"]`)
    )
      return;
    const w = window as Window & {
      dataLayer?: unknown[][];
      gtag?: (...args: unknown[]) => void;
    };
    w.dataLayer = w.dataLayer ?? [];
    w.gtag = (...args) => w.dataLayer?.push(args);
    w.gtag("js", new Date());
    w.gtag("config", measurementId);
    const script = document.createElement("script");
    script.async = true;
    script.dataset.hawpAnalytics = measurementId;
    script.src = `https://www.googletagmanager.com/gtag/js?id=${measurementId}`;
    document.head.append(script);
  };
  useEffect(() => {
    const saved = localStorage.getItem("hawp-analytics-consent-v1");
    if (saved === "accepted" || saved === "declined") {
      setChoice(saved);
      if (saved === "accepted") load();
    } else setVisible(true);
    const open = () => setVisible(true);
    window.addEventListener("hawp:analytics-settings", open);
    return () => window.removeEventListener("hawp:analytics-settings", open);
  }, []);
  const choose = (value: "accepted" | "declined") => {
    if (value === "declined") {
      const w = window as Window & { gtag?: (...args: unknown[]) => void };
      w.gtag?.("consent", "update", {
        analytics_storage: "denied",
        ad_storage: "denied",
        ad_user_data: "denied",
        ad_personalization: "denied",
      });
    }
    setChoice(value);
    localStorage.setItem("hawp-analytics-consent-v1", value);
    setVisible(false);
    if (value === "accepted") load();
  };
  return (
    <>
      {visible && (
        <section
          className="analytics-consent"
          aria-label="Analytics preference"
        >
          <div className="analytics-consent__copy">
            <strong>Help improve HAWP</strong>
            <p>
              Optional Google Analytics helps us understand which pages are
              useful. It stays off unless you allow it.
            </p>
          </div>
          <div className="analytics-consent__actions">
            <button
              type="button"
              className="analytics-consent__allow"
              onClick={() => choose("accepted")}
            >
              Allow analytics
            </button>
            <button
              type="button"
              className="analytics-consent__decline"
              onClick={() => choose("declined")}
            >
              No thanks
            </button>
          </div>
        </section>
      )}
      <button
        className="analytics-settings"
        type="button"
        aria-label="Analytics settings"
        onClick={() => setVisible(true)}
      >
        {choice === "accepted"
          ? "Analytics on · settings"
          : choice === "declined"
            ? "Analytics off · settings"
            : "Analytics settings"}
      </button>
    </>
  );
}
