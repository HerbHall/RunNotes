import { vi } from "vitest";

const mockDdClient = {
  extension: {
    vm: {
      service: {
        get: vi.fn(),
        post: vi.fn(),
        put: vi.fn(),
        delete: vi.fn(),
      },
    },
  },
  docker: {
    listContainers: vi.fn().mockResolvedValue([]),
  },
  desktopUI: {
    toast: {
      success: vi.fn(),
      warning: vi.fn(),
      error: vi.fn(),
    },
  },
  host: {
    openExternal: vi.fn(),
  },
};

export const createDockerDesktopClient = vi.fn(() => mockDdClient);
