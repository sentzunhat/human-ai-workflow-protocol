import { mkdirSync, writeFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const scriptDirectory = dirname(fileURLToPath(import.meta.url))
const outputPath = resolve(
  scriptDirectory,
  '../../.hawp/work/active/098c8a5c/hawp-featured-card.svg',
)

const artwork = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="630" viewBox="0 0 1200 630" role="img" aria-labelledby="title description">
  <title id="title">HAWP: shape intent and carry it through the work</title>
  <desc id="description">A hand-drawn workflow illustration showing a human-authored intent card passing through the HAWP shape into project work and a handoff.</desc>
  <defs>
    <pattern id="grid" width="48" height="48" patternUnits="userSpaceOnUse">
      <path d="M48 0H0V48" fill="none" stroke="#2a2022" stroke-width="1"/>
    </pattern>
  </defs>
  <rect width="1200" height="630" fill="#151112"/>
  <rect width="1200" height="630" fill="url(#grid)"/>

  <g fill="none" stroke="#c58a6d" stroke-linecap="round" stroke-linejoin="round">
    <path d="M80 72 103 56 126 73 117 99 89 99Z" stroke-width="3.5"/>
    <path d="M82 75 105 59 128 76" stroke-width="1.2" opacity=".5"/>
  </g>
  <text x="148" y="91" fill="#f6ece5" font-family="Inter,system-ui,sans-serif" font-size="30" font-weight="700" letter-spacing="1.5">HAWP</text>
  <text x="80" y="158" fill="#f6ece5" font-family="Inter,system-ui,sans-serif" font-size="38" font-weight="650">Shape intent. Carry it through the work.</text>
  <text x="82" y="197" fill="#c6b4aa" font-family="Inter,system-ui,sans-serif" font-size="19">Human-AI Workflow Protocol</text>

  <!-- Human-authored intent card -->
  <g transform="rotate(-1.2 205 393)" stroke-linecap="round" stroke-linejoin="round">
    <path d="M86 274 Q88 268 96 269 L300 271 Q307 272 307 280 L304 493 Q303 500 295 500 L94 498 Q86 497 86 490Z" fill="#f6ece5" stroke="#e5b9a2" stroke-width="3"/>
    <path d="M102 284 Q193 279 287 285" fill="none" stroke="#c58a6d" stroke-width="2"/>
    <text x="112" y="325" fill="#542d31" font-family="Inter,system-ui,sans-serif" font-size="22" font-weight="700">Intent</text>
    <path d="M112 350 Q170 346 263 351 M112 370 Q189 367 274 372 M112 390 Q166 387 236 391" fill="none" stroke="#725a55" stroke-width="3"/>
    <path d="M112 428 Q174 424 270 430 M112 448 Q177 445 248 449" fill="none" stroke="#9b8178" stroke-width="2.5"/>
    <path d="M266 454 Q279 445 288 454 Q296 463 284 475 Q273 482 264 470" fill="none" stroke="#7b2633" stroke-width="2.5"/>
  </g>

  <!-- An imperfect, hand-drawn route through the structured shape -->
  <g fill="none" stroke-linecap="round" stroke-linejoin="round">
    <path d="M320 383 C367 353 392 416 438 387 S493 357 524 385" stroke="#e5b9a2" stroke-width="4"/>
    <path d="M322 389 C370 361 393 423 440 394 S494 364 525 391" stroke="#c58a6d" stroke-width="1.5" opacity=".65"/>
    <path d="M678 387 C721 355 751 418 789 385 S829 359 868 390" stroke="#e5b9a2" stroke-width="4"/>
    <path d="M679 394 C723 362 754 426 791 392 S830 367 869 397" stroke="#c58a6d" stroke-width="1.5" opacity=".65"/>
  </g>
  <path d="m432 377 13 10-15 5" fill="none" stroke="#e5b9a2" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
  <path d="m855 380 13 10-15 5" fill="none" stroke="#e5b9a2" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>

  <!-- HAWP shape -->
  <g transform="rotate(1 600 390)" stroke-linejoin="round">
    <path d="M600 284 666 330 642 406 558 406 534 330Z" fill="#251d1f" stroke="#e5b9a2" stroke-width="4" stroke-linecap="round"/>
    <path d="M603 292 659 332 637 398 563 398 541 333Z" fill="none" stroke="#c58a6d" stroke-width="1.5" opacity=".8"/>
    <text x="600" y="365" text-anchor="middle" fill="#f6ece5" font-family="Inter,system-ui,sans-serif" font-size="21" font-weight="700" letter-spacing="1">HAWP</text>
    <text x="600" y="431" text-anchor="middle" fill="#c6b4aa" font-family="Inter,system-ui,sans-serif" font-size="16">A small shared shape</text>
  </g>

  <!-- Work and handoff cards -->
  <g stroke-linecap="round" stroke-linejoin="round">
    <path d="M892 296 Q898 290 904 295 L1010 299 Q1016 300 1015 307 L1009 420 Q1008 426 1002 425 L895 421 Q889 420 890 414Z" fill="#251d1f" stroke="#c58a6d" stroke-width="3"/>
    <path d="M915 321 Q953 319 987 324 M915 340 Q960 337 990 342 M915 359 Q952 356 980 360" fill="none" stroke="#e5b9a2" stroke-width="3"/>
    <text x="950" y="394" text-anchor="middle" fill="#f6ece5" font-family="Inter,system-ui,sans-serif" font-size="18" font-weight="650">Work</text>
    <path d="M936 435 Q942 430 948 436 L1053 439 Q1059 440 1058 447 L1053 510 Q1052 516 1046 515 L941 512 Q935 511 935 505Z" fill="#302326" stroke="#e5b9a2" stroke-width="3"/>
    <path d="M961 457 Q993 455 1031 460 M961 475 Q1005 473 1033 478" fill="none" stroke="#c6b4aa" stroke-width="2.5"/>
    <text x="996" y="501" text-anchor="middle" fill="#f6ece5" font-family="Inter,system-ui,sans-serif" font-size="17" font-weight="650">Handoff</text>
  </g>

  <text x="80" y="565" fill="#c6b4aa" font-family="Inter,system-ui,sans-serif" font-size="18">Intent · Context · Constraints · Output</text>
  <path d="M80 585 Q290 581 506 585" fill="none" stroke="#7b2633" stroke-width="2" stroke-linecap="round"/>
</svg>
`

mkdirSync(dirname(outputPath), { recursive: true })
writeFileSync(outputPath, artwork, 'utf8')
console.log(`Wrote ${outputPath}`)
