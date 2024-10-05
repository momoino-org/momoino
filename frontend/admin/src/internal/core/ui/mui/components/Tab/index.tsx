'use client';

import {
  Tabs,
  Tab as MuiTab,
  TabsProps,
  TabsOwnProps,
  Box,
} from '@mui/material';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import {
  PropsWithChildren,
  ReactNode,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { frontendOrigin } from '@/internal/core/config';

interface Tab {
  id: string;
  label: string;
  content: ReactNode;
}

interface TabProps extends TabsProps {
  label: string;
  syncToUrl: boolean;
  tabs: Tab[];
}

interface TabPanelProps {
  controller: string;
  id: string;
  currentTab: string;
}

function CustomTabPanel(props: PropsWithChildren<TabPanelProps>) {
  const { children, currentTab, id, controller, ...other } = props;

  return (
    <div
      aria-labelledby={controller}
      hidden={currentTab !== controller}
      id={id}
      role="tabpanel"
      {...other}
    >
      {currentTab === controller && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

export function Tab(props: TabProps) {
  const { label, syncToUrl, tabs, ...rest } = props;

  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [currentTabName, setCurrentTabName] = useState<string>(() => {
    if (syncToUrl) {
      const selectedTabName = searchParams.get('tab');

      if (
        selectedTabName !== null &&
        tabs.find((tab) => tab.id === selectedTabName)
      ) {
        return selectedTabName;
      }
    }

    return tabs[0].id;
  });

  const updateRoute = useCallback(
    (tabName: string) => {
      const url = new URL(pathname, frontendOrigin);
      url.searchParams.set('tab', tabName);
      router.replace(url.toString());
    },
    [pathname, router],
  );

  const handleTabChange: TabsOwnProps['onChange'] = (e, tabName: string) => {
    if (syncToUrl) {
      updateRoute(tabName);
    } else {
      setCurrentTabName(tabName);
    }
  };

  useEffect(() => {
    if (!syncToUrl) {
      return;
    }

    const tabName = searchParams.get('tab');

    if (tabName === null || !tabs.map((tab) => tab.id).includes(tabName)) {
      updateRoute(tabs[0].id);
      return;
    }

    setCurrentTabName(tabName);
  }, [pathname, router, searchParams, syncToUrl, tabs, updateRoute]);

  return (
    <>
      <Tabs
        aria-label={label}
        value={currentTabName}
        onChange={handleTabChange}
        {...rest}
      >
        {tabs.map((tab) => (
          <MuiTab
            key={tab.id}
            aria-controls={`tab-panel-${tab.id}`}
            id={tab.id}
            label={tab.label}
            value={tab.id}
          />
        ))}
      </Tabs>
      {tabs.map((tab) => (
        <CustomTabPanel
          key={tab.id}
          controller={tab.id}
          currentTab={currentTabName}
          id={`tab-panel-${tab.id}`}
        >
          {tab.content}
        </CustomTabPanel>
      ))}
    </>
  );
}
