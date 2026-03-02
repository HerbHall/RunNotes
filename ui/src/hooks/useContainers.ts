import { useState, useEffect, useCallback } from "react";
import type { ContainerInfo } from "../types";
import { listContainers } from "../api/client";

export function useContainers() {
  const [containers, setContainers] = useState<ContainerInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await listContainers();
      setContainers(result);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load containers",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return { containers, loading, error, refresh };
}
