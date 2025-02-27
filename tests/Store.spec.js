import store from "@/store";
import axios from "axios";
import jwtDecode from "jwt-decode";

jest.mock("jwt-decode", () =>
  jest.fn(() => ({
    user_id: 1,
    email: "test@example.com",
    role: "User",
  }))
);

describe("Vuex Store", () => {
  beforeEach(() => {
    localStorage.clear();
    store.replaceState({ user: null, token: null });
  });

  it("commits a token via mutation", () => {
    store.commit("setToken", "dummy-token");
    expect(store.state.token).toBe("dummy-token");
    expect(localStorage.getItem("token")).toBe("dummy-token");
  });

  it("logs out correctly via mutation", () => {
    store.commit("setToken", "dummy-token");
    store.commit("logout");
    expect(store.state.token).toBeNull();
    expect(localStorage.getItem("token")).toBeNull();
  });

  it("login action commits token on success", async () => {
    jest.spyOn(axios, "post").mockResolvedValue({ data: { token: "test-token" } });
    await store.dispatch("login", { email: "test@example.com", password: "123" });
    expect(store.state.token).toBe("test-token");
    axios.post.mockRestore();
  });

  it("register action commits token on success", async () => {
    jest.spyOn(axios, "post").mockResolvedValue({ data: { token: "register-token" } });
    await store.dispatch("register", { name: "Test", email: "test@example.com", password: "123" });
    expect(store.state.token).toBe("register-token");
    axios.post.mockRestore();
  });
});

