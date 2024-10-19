import { ReactNode } from 'react';

/**
 * Interface representing the modal context.
 */
export interface ModalContext {
  /**
   * Indicates whether the modal is visible.
   */
  visible: boolean;

  /**
   * Function to close the modal.
   */
  close: () => void;
}

/**
 * Props for the `useModal` hook.
 */
export interface UseModalProps {
  /**
   * The content to be displayed inside the modal.
   */
  content: ReactNode;

  /**
   * Optional callback to be called when the modal closes.
   */
  onClose?: () => void;
}

/**
 * Result type for the `useModal` hook.
 */
export interface UseModalResult {
  /**
   * Indicates whether the modal is visible.
   */
  visible: boolean;

  /**
   * Function to open the modal.
   */
  open: () => void;

  /**
   * Function to close the modal.
   */
  close: () => void;

  /**
   * The modal content wrapped in a provider for context.
   */
  content: ReactNode;
}

/**
 * Props for the `useModalContext` hook.
 */
export interface UseModalContextProps {
  /**
   * Optional callback to be called when the modal closes.
   */
  onClose?: () => void;
}

/**
 * Result type for the `useModalContext` hook.
 */
export interface UseModalContextResult {
  /**
   * Indicates whether the modal is currently visible.
   */
  visible: boolean;

  /**
   * Function to close the modal.
   */
  close: () => void;
}
