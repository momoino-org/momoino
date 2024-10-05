import { ProviderList } from '@/internal/modules/user-management';
import { Tab } from '@/internal/core/ui/mui/components/Tab';

export default function UserMgtPage() {
  return (
    <Tab
      syncToUrl
      label="User management page"
      tabs={[
        {
          id: 'user-management',
          label: 'User management',
          content: null,
        },
        {
          id: 'provider-management',
          label: 'Provider management',
          content: <ProviderList />,
        },
      ]}
    />
  );
}
