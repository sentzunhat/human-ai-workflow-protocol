import { useEffect, useState } from "react";
import {
  currentVersion,
  loadLatestRelease,
  releaseLoadingVersion,
  type LatestRelease,
} from "./lib/content/site";
import { Header } from "./components/shared";
import { Hero } from "./components/hero";
import {
  Agents,
  AnalyticsConsent,
  Benchmarks,
  Downloads,
  Faq,
  FinalCta,
  Footer,
  Problem,
  Protocol,
  Workflow,
} from "./components/sections";

export function App() {
  const [release, setRelease] = useState<LatestRelease | null>(null);
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    loadLatestRelease()
      .then(setRelease)
      .finally(() => setLoading(false));
  }, []);
  const version =
    release?.tag_name.replace(/^v/, "") ??
    (loading ? releaseLoadingVersion : currentVersion);
  return (
    <>
      <Header />
      <main>
        <Hero version={version} loading={loading} />
        <Problem />
        <Protocol />
        <Agents version={version} loading={loading} />
        <Benchmarks version={version} loading={loading} />
        <Workflow />
        <Downloads />
        <Faq />
        <FinalCta />
      </main>
      <Footer />
      <AnalyticsConsent />
    </>
  );
}
