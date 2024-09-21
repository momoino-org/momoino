'use client';

import { useState, useCallback } from 'react';

interface UseModalReturnType {
  visible: boolean;
  show: () => void;
  close: () => void;
}

interface UseModalProps {
  /**
   * Initial state of the modal
   */
  defaultVisible: boolean;
}

export function useModal(
  props: UseModalProps = {
    defaultVisible: false,
  },
): UseModalReturnType {
  const [visible, setVisible] = useState<boolean>(props.defaultVisible);

  const show = useCallback(() => setVisible(true), []);
  const close = useCallback(() => setVisible(false), []);

  return {
    visible,
    show,
    close,
  };
}
