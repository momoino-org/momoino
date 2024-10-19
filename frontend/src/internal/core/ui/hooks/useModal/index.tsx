'use client';

import {
  useState,
  useCallback,
  createContext,
  useContext,
  useEffect,
} from 'react';
import { Dialog } from '@mui/material';
import {
  type ModalContext,
  type UseModalContextProps,
  type UseModalContextResult,
  type UseModalProps,
  type UseModalResult,
} from './types';

/**
 * Context for managing modal state.
 */
const ModalContext = createContext<ModalContext>({
  visible: false,
  close: () => {},
});

/**
 * Custom hook to manage the modal's visibility and content.
 *
 * @param props - The properties for managing the modal.
 * @returns An object containing modal visibility state, functions to open and close the modal,
 * and the modal content wrapped in a provider.
 */
export function useModal(props: UseModalProps): UseModalResult {
  const { content, onClose } = props;
  const [visible, setVisible] = useState<boolean>(false);

  const handleOpen = useCallback(() => {
    setVisible(true);
  }, []);

  const handleClose = useCallback(() => {
    setVisible(false);
  }, []);

  useEffect(() => {
    if (visible === false) {
      onClose?.();
    }
  }, [onClose, visible]);

  return {
    visible,
    open: handleOpen,
    close: handleClose,
    content: (
      <ModalContext.Provider value={{ visible, close: handleClose }}>
        <Dialog fullWidth maxWidth="md" open={visible}>
          {content}
        </Dialog>
      </ModalContext.Provider>
    ),
  };
}

/**
 * Custom hook for accessing the modal context.
 *
 * @param props - Optional parameters for the modal context.
 * @returns The current modal context including visibility state and close function.
 */
export function useModalContext(
  props: UseModalContextProps = {},
): UseModalContextResult {
  const { onClose } = props;
  const ctx = useContext(ModalContext);

  useEffect(() => {
    if (ctx.visible === false) {
      onClose?.();
    }
  }, [ctx.visible, onClose]);

  return {
    visible: ctx.visible,
    close: ctx.close,
  };
}
