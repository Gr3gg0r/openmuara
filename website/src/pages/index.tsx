import type {ReactNode} from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

const SUPPORTED = ['Stripe', 'Fawry', 'SenangPay', 'iPay88', 'Billplz', 'ToyyibPay'];
const ROADMAP = ['Adyen', 'Xendit', 'GrabPay', "Touch 'n Go", '2C2P', 'GCash / Maya'];

function Hero(): ReactNode {
  const heroShot = useBaseUrl('/img/shots/ledger-dark.png');
  return (
    <header className={styles.hero}>
      <div className={styles.wrap}>
        <span className={styles.pill}>OPEN SOURCE · LOCAL-FIRST · OFFLINE</span>
        <Heading as="h1" className={styles.heroTitle}>
          Run the providers your app integrates with,{' '}
          <span className={styles.grad}>locally</span>.
        </Heading>
        <p className={styles.heroSub}>
          OpenMuara emulates checkout, callback, and webhook flows on your own
          machine — offline, fast, and headless. No live accounts, no tunnels,
          no surprise charges.
        </p>
        <div className={styles.heroCta}>
          <Link className={styles.btnPrimary} to="/docs/quickstart">
            Get Started
          </Link>
          <a
            className={styles.btnGhost}
            href="https://github.com/Gr3gg0r/openmuara"
            target="_blank"
            rel="noopener">
            Star on GitHub
          </a>
        </div>
        <div className={styles.browser}>
          <div className={styles.browserBar}>
            <i aria-hidden />
            <i aria-hidden />
            <i aria-hidden />
            <span className={styles.browserUrl}>http://127.0.0.1:9000/_admin</span>
          </div>
          <img
            src={heroShot}
            alt="OpenMuara dashboard ledger listing emulated Stripe and Fawry transactions with paid and new statuses"
            width={1280}
            height={512}
          />
        </div>
      </div>
    </header>
  );
}

function Showcase(): ReactNode {
  const providers = useBaseUrl('/img/shots/providers-light.png');
  const config = useBaseUrl('/img/shots/provider-config-light.png');
  return (
    <section className={styles.section}>
      <div className={styles.wrap}>
        <div className={styles.secHead}>
          <span className={styles.eyebrow}>See it working</span>
          <h2>One dashboard for the whole flow</h2>
          <p>
            Enable providers, grab the base URL, and watch every request,
            callback, and webhook land in a single ledger.
          </p>
        </div>
        <div className={styles.shotGrid}>
          <figure className={styles.shot}>
            <img
              src={providers}
              alt="OpenMuara settings view showing the emulated providers — Stripe, Fawry, iPay88, SenangPay, Billplz, and ToyyibPay — each enabled"
              width={1280}
              height={1032}
              loading="lazy"
            />
            <figcaption>Enable the providers your app uses</figcaption>
          </figure>
          <figure className={styles.shot}>
            <img
              src={config}
              alt="OpenMuara provider configuration page showing the local base URL http://127.0.0.1:9000/fawry/v1 and a sample charge endpoint"
              width={1280}
              height={1324}
              loading="lazy"
            />
            <figcaption>Copy the local base URL and point your app at it</figcaption>
          </figure>
        </div>
      </div>
    </section>
  );
}

function Providers(): ReactNode {
  return (
    <section className={`${styles.section} ${styles.sectionAlt}`}>
      <div className={styles.wrap}>
        <div className={styles.secHead}>
          <span className={styles.eyebrow}>Providers</span>
          <h2>Supported today — and growing</h2>
          <p>
            Each provider is a contract-faithful plugin: request and response
            shapes, signature verification, and a simulation page to drive the
            outcome.
          </p>
        </div>
        <div className={styles.chips}>
          {SUPPORTED.map((name) => (
            <span key={name} className={styles.chip}>
              {name}
            </span>
          ))}
          {ROADMAP.map((name) => (
            <span key={name} className={`${styles.chip} ${styles.chipRoadmap}`}>
              {name}
            </span>
          ))}
        </div>
        <p className={styles.providersNote}>
          Dashed chips are on the roadmap. Providers are plugins —{' '}
          <Link to="/docs/contributing-providers">help add yours</Link>.
        </p>
      </div>
    </section>
  );
}

const STEPS = [
  {
    title: 'Start OpenMuara',
    body: 'muara init && muara start — it listens on 127.0.0.1:9000 with a dashboard at /_admin.',
  },
  {
    title: 'Point your app at it',
    body: 'Use http://127.0.0.1:9000 as the base URL instead of the live provider. No credentials, no network.',
  },
  {
    title: 'Run the flow',
    body: 'Trigger a charge, complete it on the simulation page, and watch the callback and signed webhook reach your app.',
  },
];

function HowItWorks(): ReactNode {
  const flow = useBaseUrl('/img/flow-diagram.svg');
  return (
    <section className={styles.section}>
      <div className={styles.wrap}>
        <div className={styles.secHead}>
          <span className={styles.eyebrow}>How it works</span>
          <h2>Three steps between your app and a full test run</h2>
        </div>
        <div className={styles.steps}>
          {STEPS.map((step, i) => (
            <div key={step.title} className={styles.step}>
              <div className={styles.stepNum}>{i + 1}</div>
              <h3>{step.title}</h3>
              <p>{step.body}</p>
            </div>
          ))}
        </div>
        <div className={styles.flow}>
          <img
            src={flow}
            alt="Diagram: your app points its base URL at OpenMuara on port 9000, OpenMuara emulates the provider, you complete the outcome on a simulation page, and your app receives a signed webhook while the ledger records everything"
            width={960}
            height={340}
            loading="lazy"
          />
        </div>
      </div>
    </section>
  );
}

type TermLine =
  | {kind: 'comment'; text: string}
  | {kind: 'cmd'; text: string};

const TERM: TermLine[] = [
  {kind: 'comment', text: 'start the local engine — dashboard at /_admin'},
  {kind: 'cmd', text: 'muara init --defaults'},
  {kind: 'cmd', text: 'muara start'},
  {kind: 'comment', text: 'force an outcome — no browser needed'},
  {kind: 'cmd', text: 'muara scenario success tx-123'},
  {kind: 'comment', text: 'fire the signed webhook to your app'},
  {kind: 'cmd', text: 'muara webhook replay tx-123'},
  {kind: 'comment', text: 'script it all — JSON output for CI'},
  {kind: 'cmd', text: 'muara transaction list --json'},
];

function Cli(): ReactNode {
  return (
    <section className={`${styles.section} ${styles.sectionAlt}`}>
      <div className={styles.wrap}>
        <div className={styles.cliGrid}>
          <div className={styles.cliIntro}>
            <span className={styles.eyebrow}>Headless by design</span>
            <h2>Drive the whole flow from your terminal</h2>
            <p>
              The dashboard is optional. The <code>muara</code> CLI drives every
              step, so it slots straight into scripts and CI.
            </p>
            <ul className={styles.cliList}>
              <li>Start the engine and check it with <code>doctor</code></li>
              <li>Force an outcome — success, fail, or timeout</li>
              <li>Replay signed webhooks to your app</li>
              <li>JSON output on every command for CI</li>
            </ul>
            <Link className={styles.cliLink} to="/docs/cli">
              Read the CLI reference →
            </Link>
          </div>
          <div className={styles.terminal}>
            <div className={styles.terminalBar}>
              <span className={styles.terminalGlyph} aria-hidden>
                ❯
              </span>
              <span className={styles.terminalTitle}>muara — zsh</span>
              <span className={styles.terminalTag}>headless</span>
            </div>
            <pre className={styles.terminalBody}>
              {TERM.map((line, i) =>
                line.kind === 'comment' ? (
                  <span key={i} className={styles.tComment}>
                    # {line.text}
                    {'\n'}
                  </span>
                ) : (
                  <span key={i} className={styles.tCmd}>
                    <span className={styles.tPrompt}>❯</span> {line.text}
                    {'\n'}
                  </span>
                ),
              )}
              <span className={styles.tCmd}>
                <span className={styles.tPrompt}>❯</span>{' '}
                <span className={styles.tCursor} aria-hidden />
              </span>
            </pre>
          </div>
        </div>
      </div>
    </section>
  );
}

type Feature = {
  title: string;
  body: string;
  icon: ReactNode;
};

const iconProps = {
  viewBox: '0 0 24 24',
  'aria-hidden': true,
} as const;

const FEATURES: Feature[] = [
  {
    title: 'Local-first & offline',
    body: 'The engine and ledger run on your machine. Nothing leaves localhost — no live calls, no real balances.',
    icon: (
      <svg {...iconProps}>
        <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
        <path d="m9 12 2 2 4-4" />
      </svg>
    ),
  },
  {
    title: 'Outcomes on demand',
    body: 'Complete or cancel a flow from the simulation page, and fire the resulting webhook exactly when you want it.',
    icon: (
      <svg {...iconProps}>
        <path d="M5 3l14 7-7 2 2 7-9-16z" />
        <path d="M14 14l5 5" />
      </svg>
    ),
  },
  {
    title: 'Webhook capture & replay',
    body: 'Inspect payloads and headers, verify signatures, and replay any event from the dashboard — no tunnel required.',
    icon: (
      <svg {...iconProps}>
        <path d="M21 12a9 9 0 1 1-3-6.7" />
        <path d="M21 3v5h-5" />
      </svg>
    ),
  },
  {
    title: 'One unified ledger',
    body: 'Every request, callback, and webhook in a single time-ordered feed. Search, filter, and drill into any row.',
    icon: (
      <svg {...iconProps}>
        <path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01" />
      </svg>
    ),
  },
  {
    title: 'Contract-faithful',
    body: 'Request and response shapes and signature schemes match the real sandbox, so your integration code behaves the same.',
    icon: (
      <svg {...iconProps}>
        <rect x="3" y="3" width="18" height="18" rx="2" />
        <path d="m9 12 2 2 4-4" />
      </svg>
    ),
  },
  {
    title: 'CLI + dashboard',
    body: 'Drive it headlessly from the muara CLI or watch it live in the browser. Same engine, two surfaces.',
    icon: (
      <svg {...iconProps}>
        <rect x="2" y="4" width="20" height="14" rx="2" />
        <path d="m8 21h8M12 17v4" />
      </svg>
    ),
  },
];

function Features(): ReactNode {
  return (
    <section className={`${styles.section} ${styles.sectionAlt}`}>
      <div className={styles.wrap}>
        <div className={styles.secHead}>
          <span className={styles.eyebrow}>Why OpenMuara</span>
          <h2>Built for testing the parts that actually move money</h2>
          <p>
            Provider sandboxes are limited, and you can't trigger a real webhook
            from your laptop. OpenMuara emulates both on your machine.
          </p>
        </div>
        <div className={styles.featureGrid}>
          {FEATURES.map((f) => (
            <div key={f.title} className={styles.feature}>
              <div className={styles.featureIcon}>{f.icon}</div>
              <h3>{f.title}</h3>
              <p>{f.body}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function Etymology(): ReactNode {
  return (
    <section className={styles.section}>
      <div className={styles.wrap}>
        <div className={styles.etymology}>
          <div className={styles.etymologyWord}>
            <span>muara</span> · <i>/moo-ah-rah/</i>
          </div>
          <p>
            Malay for <b>estuary</b> — where the river meets the sea. OpenMuara
            is the calm layer where your app meets the messy world of
            integrations, before either of you goes near the real thing.
          </p>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  return (
    <Layout
      title="Run the providers your app integrates with, locally"
      description="OpenMuara emulates checkout, callback, and webhook flows on your own machine — offline, fast, and headless. No live accounts, no tunnels.">
      <Hero />
      <main>
        <Showcase />
        <Providers />
        <HowItWorks />
        <Cli />
        <Features />
        <Etymology />
      </main>
    </Layout>
  );
}
