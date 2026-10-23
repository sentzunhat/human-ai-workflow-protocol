import {
  useMemo,
  useRef,
  useState,
  type KeyboardEvent,
  type PointerEvent,
} from "react";
import { githubUrl } from "../lib/content/site";
import { ReleaseVersion } from "./shared";

type PacketNode = { x: number; y: number; label: string; detail: string };
const initialNodes: PacketNode[] = [
  { x: 50, y: 10, label: "Input", detail: "request" },
  { x: 86, y: 39, label: "Context", detail: "background" },
  { x: 75, y: 88, label: "Mission", detail: "objective" },
  { x: 25, y: 88, label: "Constraints", detail: "limits" },
  { x: 14, y: 39, label: "Output", detail: "done state" },
];
export function Hero({
  version,
  loading,
}: {
  version: string;
  loading: boolean;
}) {
  const visual = useRef<HTMLDivElement>(null);
  const [nodes, setNodes] = useState(initialNodes);
  const [dragging, setDragging] = useState(-1);
  const paths = useMemo(
    () =>
      nodes.map((node, index) => {
        const dx = node.x - 50,
          dy = node.y - 50,
          distance = Math.max(Math.hypot(dx, dy), 1);
        const px = -dy / distance,
          py = dx / distance,
          sway = [2.4, -2.8, 2.1, -2.1, 2.8][index];
        return `M 50 50 C ${(50 + dx * 0.34 + px * sway).toFixed(2)} ${(50 + dy * 0.34 + py * sway).toFixed(2)}, ${(node.x - dx * 0.34 + px * sway).toFixed(2)} ${(node.y - dy * 0.34 + py * sway).toFixed(2)}, ${node.x.toFixed(2)} ${node.y.toFixed(2)}`;
      }),
    [nodes],
  );
  const movePointer = (event: PointerEvent<HTMLDivElement>) => {
    if (dragging < 0 || !visual.current) return;
    const bounds = visual.current.getBoundingClientRect();
    const x = Math.max(
      8,
      Math.min(92, ((event.clientX - bounds.left) / bounds.width) * 100),
    );
    const y = Math.max(
      8,
      Math.min(92, ((event.clientY - bounds.top) / bounds.height) * 100),
    );
    setNodes((current) =>
      current.map((node, index) =>
        index === dragging ? { ...node, x, y } : node,
      ),
    );
  };
  const keyboardMove = (
    event: KeyboardEvent<HTMLButtonElement>,
    index: number,
  ) => {
    const steps: Record<string, [number, number]> = {
      ArrowLeft: [-1, 0],
      ArrowRight: [1, 0],
      ArrowUp: [0, -1],
      ArrowDown: [0, 1],
    };
    const delta = steps[event.key];
    if (!delta) return;
    event.preventDefault();
    const step = event.shiftKey ? 5 : 2;
    setNodes((current) =>
      current.map((node, i) =>
        i === index
          ? {
              ...node,
              x: Math.max(8, Math.min(92, node.x + delta[0] * step)),
              y: Math.max(8, Math.min(92, node.y + delta[1] * step)),
            }
          : node,
      ),
    );
  };
  return (
    <section
      id="top"
      className="site-grid relative overflow-hidden border-b border-[var(--hawp-line)]"
    >
      <div className="pointer-events-none absolute -right-40 -top-52 size-[38rem] rounded-full bg-[var(--hawp-red-deep)]/16 blur-3xl" />
      <div className="mx-auto grid min-h-[710px] max-w-[1280px] items-center gap-14 px-5 py-20 sm:px-8 lg:grid-cols-[1.08fr_.92fr] lg:px-10 lg:py-24">
        <div className="relative z-10">
          <div className="mb-8 flex items-center gap-3 text-[11px] font-semibold tracking-[0.08em] text-[var(--hawp-peach)] uppercase">
            <span className="size-3 rounded-full bg-[var(--hawp-red)]" />
            <ReleaseVersion version={version} loading={loading} badge /> · Open
            protocol · Markdown-first · MCP-ready
          </div>
          <h1 className="max-w-[760px] text-5xl font-bold leading-[1.02] tracking-[-0.045em] text-[var(--hawp-ink)] sm:text-6xl lg:text-[5rem]">
            Shape the work before the model{" "}
            <span className="text-[var(--hawp-peach)]">touches it.</span>
          </h1>
          <p className="mt-7 max-w-[730px] text-lg leading-7 text-[var(--hawp-muted)] sm:text-xl sm:leading-8">
            HAWP gives AI agent work a stable shape for intent, context,
            mission, constraints, and output, so agents can execute with less
            drift and hand off without starting from zero.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <a
              href={githubUrl}
              className="rounded-xl border border-[var(--hawp-red)] bg-[var(--hawp-red-deep)] px-6 py-3.5 text-sm font-semibold text-white transition hover:bg-[var(--hawp-red)]"
            >
              Explore on GitHub ↗
            </a>
            <a
              href="#protocol"
              className="rounded-xl border border-[var(--hawp-line-strong)] bg-[var(--hawp-panel-soft)] px-6 py-3.5 text-sm font-semibold text-[var(--hawp-ink)]"
            >
              See the protocol ↓
            </a>
          </div>
          <p className="mt-6 text-xs text-[var(--hawp-muted-2)]">
            Apache-2.0 · Go CLI · Static Markdown artifacts · No lock-in
          </p>
        </div>
        <div
          ref={visual}
          className="protocol-visual"
          role="group"
          aria-label="Five HAWP fields converging into a stable work packet"
          onPointerMove={movePointer}
          onPointerUp={() => setDragging(-1)}
          onPointerCancel={() => setDragging(-1)}
        >
          <div className="orbit one" />
          <div className="orbit two" />
          <div className="orbit three" />
          <svg className="connections" viewBox="0 0 100 100" aria-hidden="true">
            {paths.map((path, index) => (
              <g key={index}>
                <path
                  className="connector-shadow"
                  d={path}
                  pathLength="1"
                  style={
                    {
                      "--connector-delay": `${index * -0.35}s`,
                    } as React.CSSProperties
                  }
                />
                <path className="connector-base" d={path} />
                <path
                  className="connector"
                  d={path}
                  pathLength="1"
                  style={
                    {
                      "--connector-delay": `${index * -0.35}s`,
                    } as React.CSSProperties
                  }
                />
              </g>
            ))}
          </svg>
          <div className="core">
            <svg viewBox="0 0 100 100" aria-hidden="true">
              <path d="M50 7 90.9 36.7 75.3 84.8 24.7 84.8 9.1 36.7Z" />
            </svg>
            <span>WORK</span>
          </div>
          {nodes.map((node, index) => (
            <button
              key={node.label}
              type="button"
              className={`node${dragging === index ? " dragging" : ""}`}
              style={
                {
                  "--node-x": `${node.x}%`,
                  "--node-y": `${node.y}%`,
                  "--float-delay": `${index * -0.7}s`,
                } as React.CSSProperties
              }
              aria-describedby="drag-instruction"
              onPointerDown={(event) => {
                if (event.pointerType === "mouse" && event.button !== 0) return;
                setDragging(index);
                event.currentTarget.setPointerCapture(event.pointerId);
                event.preventDefault();
              }}
              onPointerUp={() => setDragging(-1)}
              onKeyDown={(event) => keyboardMove(event, index)}
            >
              <span className="node-content">
                <strong>{node.label}</strong>
                {node.detail}
              </span>
            </button>
          ))}
          <p id="drag-instruction" className="drag-hint">
            Drag a field or focus it and use arrow keys to reshape the work
            packet
          </p>
        </div>
      </div>
    </section>
  );
}
