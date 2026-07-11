# Benchmark Comparison Report

## Overview
This report compares two different approaches to responding to a vague, high-level task: an unguided response (No HAWP) and a structured response using the Human-AI Workflow Protocol (HAWP).

## Methodology
**Task:** "review this codebase and tell me what's wrong"

1.  **Approach 1: No HAWP (Unguided)**
    *   The AI was given the task without any framing, context, or constraints.
    *   The response was a general overview of the repository structure and potential issues based on surface-level observations.
2.  **Approach 2: With HAWP (Structured)**
    *   The AI used the HAWP "Shape" pattern to define the task's scope before responding.
    *   This involved explicitly defining `input`, `context`, `mission`, `constraints`, and `output` formats.
    *   The response was a highly targeted, evidence-based analysis of specific risks (Complexity, Documentation Fragmentation, Maintenance Scalability) with exact file paths as evidence.

## Comparison Analysis

| Feature | No HAWP (Unguided) | With HAWP (Structured) |
| :--- | :--- | :--- |
| **Scope Control** | Low. The response was broad and covered everything from structure to potential issues without clear boundaries. | High. The `constraints` field explicitly limited the review to top-level directories and `.github/instructions`. |
| **Precision & Evidence** | Low. Observations were general (e.g., "The presence of many different folders... indicates a high level of complexity"). No specific file paths or evidence provided. | High. Every identified risk was backed by specific, repo-relative file paths as evidence (e.g., `librarian/scripts/backlog-upgrade/index.ts`). |
| **Actionability** | Low. The findings were vague and could be interpreted in many ways. They provided a "summary" but no clear starting point for investigation. | High. The findings were highly actionable, providing specific locations of potential risk (e.s., `distribution/sources/providers/continue/` vs `cursor`). |
| **Risk Identification** | Surface-scale. Identified general architectural patterns and potential issues without deep analysis. | Deep & Targeted. Focused on identifying high-risk areas that could impede developer onboarding or maintainability. |
| **Structure of Output** | Unstructured text with bullet points. | Highly structured, following a predefined `output` schema defined in the HAWP shape. |

## Conclusion
The use of HAWP significantly improves the quality of AI-generated responses for complex tasks. By forcing the AI to "shape" the task first—defining its own boundaries, mission, and constraints—the response becomes much more precise, evidence-based, and actionable. The structured approach reduces "hallucination" or speculative findings by explicitly stating what is observable from the file tree and READMEs.
