import type { ReactNode } from "react";
import clsx from "clsx";
import Heading from "@theme/Heading";
import styles from "./styles.module.css";

type FeatureItem = {
  title: string;
  Svg?: React.ComponentType<React.ComponentProps<"svg">>;
  description: ReactNode;
};

const FeatureList: FeatureItem[] = [
  {
    title: "Minimalist & Essential",
    description: (
      <>
        Backupman is a compact open-source solution focusing on the essentials
        of database backups: creating dumps, store dump, and automation.
      </>
    ),
  },
  {
    title: "Automated & Reliable",
    description: (
      <>
        Set up scheduled backups with intelligent retention rules, ensuring your
        data is regularly secured and old backups are automatically cleaned up.
      </>
    ),
  },
  {
    title: "Flexible Storage & Alerts",
    description: (
      <>
        Store your backups locally or in cloud services like Google Drive, and
        stay informed with real-time notifications about backup status.
      </>
    ),
  },
];

function Feature({ title, description }: FeatureItem) {
  return (
    <div className={clsx("col col--4")}>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
