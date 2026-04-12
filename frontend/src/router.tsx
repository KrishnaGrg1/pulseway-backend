import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import * as TanstackQuery from "./integrations/tanstack-query /root-provider";
export function getRouter() {
  const rqContext = TanstackQuery.getContext();
  const router = createTanStackRouter({
    routeTree,
    ...rqContext,
    scrollRestoration: true,
    defaultPreload: "intent",
    defaultPreloadStaleTime: 0,
  });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof getRouter>;
  }
}
